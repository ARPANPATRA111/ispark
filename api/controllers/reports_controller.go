package controllers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iips-oss/ispark/api/config"
	"github.com/iips-oss/ispark/api/models"
	"gorm.io/gorm"
)

// reportUploadDir is where generated report files are stored. It sits under the
// same uploads directory that certificates use, which is the persisted volume
// in the container.
const reportUploadDir = "./uploads/reports"

// allowedReportTypes is the set of report types the Reports Center can generate.
var allowedReportTypes = map[string]bool{
	"Student Performance":      true,
	"Mentor Analytics":         true,
	"Activity Participation":   true,
	"Credit Distribution":      true,
	"Department Analytics":     true,
	"Certificate Verification": true,
	"Leaderboard Rankings":     true,
	"Semester Summary":         true,
}

// allowedReportFormats is the set of output formats a report may request. Only
// CSV is supported for now — the backend produces CSV files, so accepting PDF or
// Excel would hand back a mislabelled file. Real PDF/Excel rendering is a future
// enhancement; until then the API and UI expose CSV only.
var allowedReportFormats = map[string]bool{
	"CSV": true,
}

// allowedExportTypes is the set of direct data exports the Export Center offers.
var allowedExportTypes = map[string]bool{
	"students":     true,
	"activities":   true,
	"credits":      true,
	"certificates": true,
}

// reportTemplate is one of the Reports Center's quick-start templates.
type reportTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// reportTemplates are the fixed quick templates shown in the Reports Center.
var reportTemplates = []reportTemplate{
	{Name: "Student Credit Summary", Description: "Credits per student", Type: "Credit Distribution"},
	{Name: "Certificate Verification", Description: "Approval status log", Type: "Certificate Verification"},
	{Name: "Activity Participation", Description: "Enrollment vs completion", Type: "Activity Participation"},
	{Name: "Batch Performance", Description: "Batch-wise comparison", Type: "Department Analytics"},
	{Name: "Mentor Activity", Description: "Mentor engagement", Type: "Mentor Analytics"},
	{Name: "Department Performance", Description: "Dept-wise metrics", Type: "Department Analytics"},
	{Name: "Semester Summary", Description: "Per-semester rollup", Type: "Semester Summary"},
	{Name: "System Usage Report", Description: "Platform activity", Type: "Student Performance"},
}

var semesterDigits = regexp.MustCompile(`\d+`)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// currentAdminName resolves a human-friendly name for the acting super admin,
// used to attribute generated reports and audit-log entries.
func currentAdminName(c *fiber.Ctx) string {
	id, _ := c.Locals("roll_no").(string)
	if id == "" {
		return "Super Admin"
	}

	var admin models.Admin
	if err := config.DB.Select("name").Where("admin_id = ?", id).First(&admin).Error; err == nil && admin.Name != "" {
		return admin.Name
	}

	return id
}

// logReportActivity records a report-related action in the audit log. It is
// best-effort: a failure to log must never fail the request that triggered it.
func logReportActivity(category, action, format, user string) {
	_ = config.DB.Create(&models.ReportAuditLog{
		Action:   action,
		Category: category,
		Format:   format,
		User:     user,
	}).Error
}

// paginate reads limit/offset query params, applying sane defaults and caps.
func paginate(c *fiber.Ctx) (limit, offset int) {
	limit = 50
	if v, err := strconv.Atoi(c.Query("limit")); err == nil && v > 0 {
		limit = v
	}
	if limit > 200 {
		limit = 200
	}

	if v, err := strconv.Atoi(c.Query("offset")); err == nil && v > 0 {
		offset = v
	}

	return limit, offset
}

// filterCourse returns the course filter value, or "" when no real filter is set.
func filterCourse(course string) string {
	course = strings.TrimSpace(course)
	if course == "" || strings.EqualFold(course, "All Courses") {
		return ""
	}
	return course
}

// filterSemester extracts a numeric semester from the UI's label (e.g.
// "Semester 2" -> 2). It returns ok=false when no specific semester is selected.
func filterSemester(semester string) (int, bool) {
	semester = strings.TrimSpace(semester)
	if semester == "" || strings.EqualFold(semester, "All Semesters") {
		return 0, false
	}
	if match := semesterDigits.FindString(semester); match != "" {
		if n, err := strconv.Atoi(match); err == nil {
			return n, true
		}
	}
	return 0, false
}

// parseReportDate parses an optional YYYY-MM-DD date from the generate form.
func parseReportDate(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// nextDay returns midnight of the day after t. It is the exclusive upper bound
// for an inclusive "to" date: a timestamp comparison of `< nextDay(to)` includes
// every record on the selected day, whereas `<= to` (which is midnight of that
// day) would drop everything after 00:00:00 and, for a single-day range, return
// nothing at all.
func nextDay(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

// ---------------------------------------------------------------------------
// Report data aggregation
// ---------------------------------------------------------------------------

// studentAggregate is a per-student rollup used by every student-centred report.
type studentAggregate struct {
	RollNo        string
	Name          string
	CourseName    string
	Semester      int
	Credits       int
	ActivityCount int
	Pending       int
}

// buildStudentAggregates loads per-student credit, activity and pending-cert
// figures, honouring the course/semester filters. Certificates and enrollments
// are aggregated in separate queries and merged in memory so a join never
// multiplies the credit sums.
func buildStudentAggregates(input models.GenerateReportInput) ([]studentAggregate, error) {
	studentQuery := config.DB.Model(&models.Student{}).
		Select("roll_no", "name", "course_name", "semester").
		Order("course_name asc, name asc")

	if course := filterCourse(input.Course); course != "" {
		studentQuery = studentQuery.Where("course_name = ?", course)
	}
	if sem, ok := filterSemester(input.Semester); ok {
		studentQuery = studentQuery.Where("semester = ?", sem)
	}

	var students []models.Student
	if err := studentQuery.Find(&students).Error; err != nil {
		return nil, err
	}

	// The date range windows the underlying data so the filter is honoured
	// consistently: credits and pending counts come from certificates whose
	// activity date falls in the range, and activity counts from enrollments
	// created within it. This matches how the certificate and activity reports
	// apply the same filter. A parse error leaves the bound unset (nil).
	from, _ := parseReportDate(input.DateFrom)
	to, _ := parseReportDate(input.DateTo)

	// When a date range is selected it also determines which students appear:
	// only those with reportable activity in the window are included. Without a
	// range, every student (matching the course/semester filters) is reported.
	dateRangeSet := from != nil || to != nil

	type certRow struct {
		StudentRollNo string
		Credits       int
		Pending       int
	}
	certQuery := config.DB.Model(&models.Certificate{}).
		Select("student_roll_no, " +
			"COALESCE(SUM(CASE WHEN status = 'Approved' THEN credits ELSE 0 END), 0) AS credits, " +
			"COUNT(CASE WHEN status = 'Pending' THEN 1 END) AS pending").
		Group("student_roll_no")
	if from != nil {
		certQuery = certQuery.Where("activity_date >= ?", *from)
	}
	if to != nil {
		certQuery = certQuery.Where("activity_date < ?", nextDay(*to))
	}
	var certRows []certRow
	if err := certQuery.Scan(&certRows).Error; err != nil {
		return nil, err
	}

	type enrRow struct {
		StudentRollNo string
		Count         int
	}
	enrQuery := config.DB.Model(&models.Enrollment{}).
		Select("student_roll_no, COUNT(*) AS count").
		Group("student_roll_no")
	if from != nil {
		enrQuery = enrQuery.Where("created_at >= ?", *from)
	}
	if to != nil {
		enrQuery = enrQuery.Where("created_at < ?", nextDay(*to))
	}
	var enrRows []enrRow
	if err := enrQuery.Scan(&enrRows).Error; err != nil {
		return nil, err
	}

	credits := make(map[string]certRow, len(certRows))
	for _, row := range certRows {
		credits[row.StudentRollNo] = row
	}
	activities := make(map[string]int, len(enrRows))
	for _, row := range enrRows {
		activities[row.StudentRollNo] = row.Count
	}

	aggregates := make([]studentAggregate, 0, len(students))
	for _, student := range students {
		cert := credits[student.RollNo]
		activityCount := activities[student.RollNo]

		// A date-ranged report only lists students who actually have reportable
		// activity in that window — earned credits, enrollments, or pending
		// certificates. This keeps a future/empty range from returning every
		// student as a row of zeros.
		if dateRangeSet && cert.Credits == 0 && activityCount == 0 && cert.Pending == 0 {
			continue
		}

		aggregates = append(aggregates, studentAggregate{
			RollNo:        student.RollNo,
			Name:          student.Name,
			CourseName:    student.CourseName,
			Semester:      student.Semester,
			Credits:       cert.Credits,
			ActivityCount: activityCount,
			Pending:       cert.Pending,
		})
	}

	return aggregates, nil
}

// buildReportData produces the CSV rows (header first) for a report type,
// applying the report's filters. It is the single place that maps a report type
// to a query.
func buildReportData(input models.GenerateReportInput) ([][]string, error) {
	switch input.Type {
	case "Student Performance":
		aggregates, err := buildStudentAggregates(input)
		if err != nil {
			return nil, err
		}
		rows := [][]string{{"Roll No", "Name", "Course", "Semester", "Credits Earned", "Activities", "Pending Certificates"}}
		for _, a := range aggregates {
			rows = append(rows, []string{a.RollNo, a.Name, a.CourseName, strconv.Itoa(a.Semester),
				strconv.Itoa(a.Credits), strconv.Itoa(a.ActivityCount), strconv.Itoa(a.Pending)})
		}
		return rows, nil

	case "Credit Distribution":
		aggregates, err := buildStudentAggregates(input)
		if err != nil {
			return nil, err
		}
		rows := [][]string{{"Roll No", "Name", "Course", "Semester", "Credits Earned"}}
		for _, a := range aggregates {
			rows = append(rows, []string{a.RollNo, a.Name, a.CourseName, strconv.Itoa(a.Semester), strconv.Itoa(a.Credits)})
		}
		return rows, nil

	case "Leaderboard Rankings":
		aggregates, err := buildStudentAggregates(input)
		if err != nil {
			return nil, err
		}
		// Highest credits first; that ordering is the ranking.
		sortByCreditsDesc(aggregates)
		rows := [][]string{{"Rank", "Roll No", "Name", "Course", "Credits Earned", "Activities"}}
		for i, a := range aggregates {
			rows = append(rows, []string{strconv.Itoa(i + 1), a.RollNo, a.Name, a.CourseName,
				strconv.Itoa(a.Credits), strconv.Itoa(a.ActivityCount)})
		}
		return rows, nil

	case "Department Analytics":
		aggregates, err := buildStudentAggregates(input)
		if err != nil {
			return nil, err
		}
		return departmentRollup(aggregates), nil

	case "Semester Summary":
		aggregates, err := buildStudentAggregates(input)
		if err != nil {
			return nil, err
		}
		return semesterRollup(aggregates), nil

	case "Certificate Verification":
		return certificateReportData(input)

	case "Activity Participation":
		return activityReportData(input)

	case "Mentor Analytics":
		return mentorReportData()

	default:
		return nil, fmt.Errorf("unsupported report type %q", input.Type)
	}
}

// sortByCreditsDesc orders aggregates by credits earned, highest first, using a
// simple insertion sort to avoid pulling in sort just for this.
func sortByCreditsDesc(aggregates []studentAggregate) {
	for i := 1; i < len(aggregates); i++ {
		current := aggregates[i]
		j := i - 1
		for j >= 0 && aggregates[j].Credits < current.Credits {
			aggregates[j+1] = aggregates[j]
			j--
		}
		aggregates[j+1] = current
	}
}

// departmentRollup groups per-student aggregates into per-course figures.
func departmentRollup(aggregates []studentAggregate) [][]string {
	type dept struct {
		students   int
		credits    int
		activities int
		order      int
	}
	depts := map[string]*dept{}
	var names []string
	for _, a := range aggregates {
		d, ok := depts[a.CourseName]
		if !ok {
			d = &dept{order: len(names)}
			depts[a.CourseName] = d
			names = append(names, a.CourseName)
		}
		d.students++
		d.credits += a.Credits
		d.activities += a.ActivityCount
	}

	rows := [][]string{{"Department", "Students", "Total Credits", "Avg Credits", "Activities"}}
	for _, name := range names {
		d := depts[name]
		avg := 0.0
		if d.students > 0 {
			avg = float64(d.credits) / float64(d.students)
		}
		rows = append(rows, []string{name, strconv.Itoa(d.students), strconv.Itoa(d.credits),
			fmt.Sprintf("%.1f", avg), strconv.Itoa(d.activities)})
	}
	return rows
}

// semesterRollup groups per-student aggregates into per-semester figures.
func semesterRollup(aggregates []studentAggregate) [][]string {
	type sem struct {
		students int
		credits  int
		pending  int
	}
	sems := map[int]*sem{}
	var order []int
	for _, a := range aggregates {
		s, ok := sems[a.Semester]
		if !ok {
			s = &sem{}
			sems[a.Semester] = s
			order = append(order, a.Semester)
		}
		s.students++
		s.credits += a.Credits
		s.pending += a.Pending
	}

	// Order semesters ascending for a stable, readable rollup.
	for i := 1; i < len(order); i++ {
		v := order[i]
		j := i - 1
		for j >= 0 && order[j] > v {
			order[j+1] = order[j]
			j--
		}
		order[j+1] = v
	}

	rows := [][]string{{"Semester", "Students", "Total Credits", "Avg Credits", "Pending Certificates"}}
	for _, key := range order {
		s := sems[key]
		avg := 0.0
		if s.students > 0 {
			avg = float64(s.credits) / float64(s.students)
		}
		rows = append(rows, []string{strconv.Itoa(key), strconv.Itoa(s.students), strconv.Itoa(s.credits),
			fmt.Sprintf("%.1f", avg), strconv.Itoa(s.pending)})
	}
	return rows
}

// certificateReportData lists certificates, honouring the course filter (via the
// owning student) and the date range (on the certificate's activity date).
func certificateReportData(input models.GenerateReportInput) ([][]string, error) {
	query := config.DB.Model(&models.Certificate{}).
		Select("certificates.student_roll_no, certificates.activity_name, certificates.activity_category, " +
			"certificates.status, certificates.credits, certificates.activity_date").
		Order("certificates.created_at desc")

	// Course and semester both live on the owning student, so a single join
	// backs either filter. Joining only when a filter is set keeps the plain
	// (unfiltered) query cheap.
	course := filterCourse(input.Course)
	sem, hasSem := filterSemester(input.Semester)
	if course != "" || hasSem {
		query = query.Joins("JOIN students ON students.roll_no = certificates.student_roll_no")
		if course != "" {
			query = query.Where("students.course_name = ?", course)
		}
		if hasSem {
			query = query.Where("students.semester = ?", sem)
		}
	}
	if from, err := parseReportDate(input.DateFrom); err == nil && from != nil {
		query = query.Where("certificates.activity_date >= ?", *from)
	}
	if to, err := parseReportDate(input.DateTo); err == nil && to != nil {
		query = query.Where("certificates.activity_date < ?", nextDay(*to))
	}

	type row struct {
		StudentRollNo    string
		ActivityName     string
		ActivityCategory string
		Status           string
		Credits          int
		ActivityDate     time.Time
	}
	var results []row
	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	rows := [][]string{{"Student", "Activity", "Category", "Status", "Credits", "Activity Date"}}
	for _, r := range results {
		rows = append(rows, []string{r.StudentRollNo, r.ActivityName, r.ActivityCategory, r.Status,
			strconv.Itoa(r.Credits), r.ActivityDate.Format("2006-01-02")})
	}
	return rows, nil
}

// activityReportData lists activities with their enrollment and completion counts.
// The course/semester filters scope the counts to enrolled students in that
// course/semester: the student attributes are matched in the join's ON clause
// (not a WHERE) so an activity with no matching enrolments still appears with a
// zero count rather than dropping out of the report entirely.
func activityReportData(input models.GenerateReportInput) ([][]string, error) {
	// The counts are over students.roll_no rather than enrollments.id so they
	// only include enrolments whose student survived the course/semester filter
	// in the students join.
	query := config.DB.Model(&models.Activity{}).
		Select("activities.name, activities.category, activities.credits, " +
			"COUNT(students.roll_no) AS enrolled, " +
			"COUNT(CASE WHEN enrollments.status = 'Completed' AND students.roll_no IS NOT NULL THEN 1 END) AS completed").
		Joins("LEFT JOIN enrollments ON enrollments.activity_id = activities.id AND enrollments.deleted_at IS NULL")

	studentJoin := "LEFT JOIN students ON students.roll_no = enrollments.student_roll_no AND students.deleted_at IS NULL"
	var studentArgs []interface{}
	if course := filterCourse(input.Course); course != "" {
		studentJoin += " AND students.course_name = ?"
		studentArgs = append(studentArgs, course)
	}
	if sem, ok := filterSemester(input.Semester); ok {
		studentJoin += " AND students.semester = ?"
		studentArgs = append(studentArgs, sem)
	}
	query = query.Joins(studentJoin, studentArgs...).
		Group("activities.id, activities.name, activities.category, activities.credits").
		Order("enrolled desc")

	if from, err := parseReportDate(input.DateFrom); err == nil && from != nil {
		query = query.Where("activities.activity_date >= ?", *from)
	}
	if to, err := parseReportDate(input.DateTo); err == nil && to != nil {
		query = query.Where("activities.activity_date < ?", nextDay(*to))
	}

	type row struct {
		Name      string
		Category  string
		Credits   int
		Enrolled  int
		Completed int
	}
	var results []row
	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	rows := [][]string{{"Activity", "Category", "Credits", "Enrolled", "Completed"}}
	for _, r := range results {
		rows = append(rows, []string{r.Name, r.Category, strconv.Itoa(r.Credits),
			strconv.Itoa(r.Enrolled), strconv.Itoa(r.Completed)})
	}
	return rows, nil
}

// mentorReportData lists admin accounts and how many students fall under each
// one's assigned batch.
func mentorReportData() ([][]string, error) {
	var admins []models.Admin
	if err := config.DB.Select("admin_id", "name", "role", "assigned_batch").
		Order("name asc").Find(&admins).Error; err != nil {
		return nil, err
	}

	rows := [][]string{{"Admin ID", "Name", "Role", "Assigned Batch", "Students"}}
	for _, admin := range admins {
		batch := admin.AssignedBatch

		// AssignedBatch is a roll-number prefix (e.g. "IT2K24"), not a course
		// name, so students are scoped by `roll_no LIKE prefix%` — the same rule
		// the admin dashboard uses. An empty batch means the account is not scoped
		// to a batch (super admin), so it counts every student.
		studentQuery := config.DB.Model(&models.Student{})
		if batch != "" {
			studentQuery = studentQuery.Where("roll_no LIKE ?", batch+"%")
		} else {
			batch = "All Batches"
		}

		var students int64
		if err := studentQuery.Count(&students).Error; err != nil {
			return nil, err
		}
		rows = append(rows, []string{admin.AdminID, admin.Name, admin.Role, batch, strconv.FormatInt(students, 10)})
	}
	return rows, nil
}

// writeReportFile writes CSV rows to a new file under the reports directory and
// returns the file name, full path and size.
func writeReportFile(reportType string, rows [][]string) (fileName, filePath string, size int64, err error) {
	if err = os.MkdirAll(reportUploadDir, 0o755); err != nil {
		return "", "", 0, err
	}

	slug := strings.ToLower(strings.ReplaceAll(reportType, " ", "_"))
	fileName = fmt.Sprintf("%s_%d.csv", slug, time.Now().UnixNano())
	filePath = filepath.Join(reportUploadDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", "", 0, err
	}
	// Cleanup on any early return; the success path closes explicitly below and
	// checks the error, since a failed close means the report was not fully written.
	defer func() { _ = file.Close() }()

	writer := csv.NewWriter(file)
	if err = writer.WriteAll(rows); err != nil {
		return "", "", 0, err
	}
	writer.Flush()
	if err = writer.Error(); err != nil {
		return "", "", 0, err
	}

	info, err := file.Stat()
	if err != nil {
		return "", "", 0, err
	}

	if err = file.Close(); err != nil {
		return "", "", 0, err
	}

	return fileName, filePath, info.Size(), nil
}

// ---------------------------------------------------------------------------
// Overview
// ---------------------------------------------------------------------------

// GetReportsSummary returns the four Reports Center overview figures.
func GetReportsSummary(c *fiber.Ctx) error {
	var (
		totalReports     int64
		scheduledReports int64
		monthlyDownloads int64
		storageBytes     int64
	)

	if err := config.DB.Model(&models.GeneratedReport{}).Count(&totalReports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load report summary"})
	}
	if err := config.DB.Model(&models.ScheduledReport{}).Count(&scheduledReports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load report summary"})
	}

	monthStart := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, -time.Now().UTC().Day()+1)
	if err := config.DB.Model(&models.ReportAuditLog{}).
		Where("category = ? AND created_at >= ?", "Download", monthStart).
		Count(&monthlyDownloads).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load report summary"})
	}

	if err := config.DB.Model(&models.GeneratedReport{}).
		Select("COALESCE(SUM(file_size), 0)").Scan(&storageBytes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load report summary"})
	}

	return c.JSON(fiber.Map{
		"total_reports":     totalReports,
		"scheduled_reports": scheduledReports,
		"monthly_downloads": monthlyDownloads,
		"storage_bytes":     storageBytes,
	})
}

// maxProgramSemester is the highest semester number a report can be filtered by.
// The semester filter offers the full 1..N range rather than only the semesters
// that currently have students, so any semester can be filtered even when empty.
const maxProgramSemester = 10

// GetReportFilters returns the course and semester options for the generate-report
// form. Both come from canonical domains, not from distinct row values: courses
// are the platform's fixed program list (models.CanonicalCourses) and semesters
// are the full 1..maxProgramSemester range. This keeps the dropdowns stable and
// consistent, so they never drift with legacy/registered course-name variants.
func GetReportFilters(c *fiber.Ctx) error {
	semesters := make([]int, 0, maxProgramSemester)
	for s := 1; s <= maxProgramSemester; s++ {
		semesters = append(semesters, s)
	}

	return c.JSON(fiber.Map{"courses": models.CanonicalCourses, "semesters": semesters})
}

// ---------------------------------------------------------------------------
// Report generation & management
// ---------------------------------------------------------------------------

// GenerateReport builds a new report from live platform data, writes it to a
// file and records it. Generation runs inline: the report row is created in the
// Processing state and transitions to Ready once its file is written, or Failed
// if aggregation or writing fails.
func GenerateReport(c *fiber.Ctx) error {
	var input models.GenerateReportInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	input.Type = strings.TrimSpace(input.Type)
	input.Format = strings.TrimSpace(input.Format)
	if input.Format == "" {
		input.Format = "CSV"
	}

	if !allowedReportTypes[input.Type] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unknown report type"})
	}
	if !allowedReportFormats[input.Format] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format must be CSV"})
	}

	dateFrom, err := parseReportDate(input.DateFrom)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid from date, expected YYYY-MM-DD"})
	}
	dateTo, err := parseReportDate(input.DateTo)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid to date, expected YYYY-MM-DD"})
	}
	if dateFrom != nil && dateTo != nil && dateTo.Before(*dateFrom) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "The to date cannot be before the from date"})
	}

	generatedBy := currentAdminName(c)

	report := models.GeneratedReport{
		Name:        input.Type + " Report",
		Type:        input.Type,
		Course:      strings.TrimSpace(input.Course),
		Semester:    strings.TrimSpace(input.Semester),
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		Format:      input.Format,
		Status:      "Processing",
		GeneratedBy: generatedBy,
	}
	if err := config.DB.Create(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create report"})
	}

	// markFailed flips the report to Failed and returns a 500 to the caller.
	markFailed := func() error {
		config.DB.Model(&report).Update("status", "Failed")
		report.Status = "Failed"
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate report"})
	}

	rows, err := buildReportData(input)
	if err != nil {
		return markFailed()
	}

	fileName, filePath, size, err := writeReportFile(input.Type, rows)
	if err != nil {
		return markFailed()
	}

	report.Status = "Ready"
	report.FileName = fileName
	report.FilePath = filePath
	report.FileSize = size
	if err := config.DB.Save(&report).Error; err != nil {
		return markFailed()
	}

	logReportActivity("Generate", report.Name+" generated", input.Format, generatedBy)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"report": report})
}

// GetGeneratedReports returns generated reports newest first, with pagination.
func GetGeneratedReports(c *fiber.Ctx) error {
	limit, offset := paginate(c)

	query := config.DB.Model(&models.GeneratedReport{})
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load reports"})
	}

	var reports []models.GeneratedReport
	if err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&reports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load reports"})
	}

	return c.JSON(fiber.Map{"reports": reports, "total": total, "limit": limit, "offset": offset})
}

// GetReportDetail returns a single report's metadata and status, which the UI
// polls while a report is still Processing.
func GetReportDetail(c *fiber.Ctx) error {
	var report models.GeneratedReport
	if err := config.DB.Where("id = ?", c.Params("id")).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Report not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load report"})
	}

	return c.JSON(fiber.Map{"report": report})
}

// DownloadReport streams a generated report's file.
func DownloadReport(c *fiber.Ctx) error {
	var report models.GeneratedReport
	if err := config.DB.Where("id = ?", c.Params("id")).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Report not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load report"})
	}

	if report.Status != "Ready" || report.FilePath == "" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Report is not ready for download"})
	}

	if _, err := os.Stat(report.FilePath); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Report file is no longer available"})
	}

	logReportActivity("Download", report.Name+" downloaded", report.Format, currentAdminName(c))

	return c.Download(report.FilePath, report.FileName)
}

// DeleteReport removes a generated report and its file.
func DeleteReport(c *fiber.Ctx) error {
	var report models.GeneratedReport
	if err := config.DB.Where("id = ?", c.Params("id")).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Report not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load report"})
	}

	if report.FilePath != "" {
		// A missing file must not block deleting the record.
		_ = os.Remove(report.FilePath)
	}

	if err := config.DB.Delete(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete report"})
	}

	return c.JSON(fiber.Map{"message": "Report deleted successfully"})
}

// GetReportTemplates returns the Reports Center's quick-start templates.
func GetReportTemplates(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"templates": reportTemplates})
}

// ---------------------------------------------------------------------------
// Scheduled reports
// ---------------------------------------------------------------------------

// GetScheduledReports returns every scheduled report.
func GetScheduledReports(c *fiber.Ctx) error {
	var scheduled []models.ScheduledReport
	if err := config.DB.Order("created_at desc").Find(&scheduled).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load scheduled reports"})
	}

	return c.JSON(fiber.Map{"scheduled": scheduled})
}

// CreateScheduledReport registers a new recurring report.
func CreateScheduledReport(c *fiber.Ctx) error {
	var input models.ScheduledReportInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Type = strings.TrimSpace(input.Type)
	input.Frequency = strings.TrimSpace(input.Frequency)
	input.Format = strings.TrimSpace(input.Format)
	if input.Format == "" {
		input.Format = "CSV"
	}

	if input.Name == "" || input.Frequency == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name and frequency are required"})
	}
	if !allowedReportTypes[input.Type] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unknown report type"})
	}
	if !allowedReportFormats[input.Format] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format must be CSV"})
	}

	enabled := true
	if input.Enabled != nil {
		enabled = *input.Enabled
	}

	scheduled := models.ScheduledReport{
		Name:      input.Name,
		Type:      input.Type,
		Frequency: input.Frequency,
		Format:    input.Format,
		Enabled:   enabled,
	}
	if err := config.DB.Create(&scheduled).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create scheduled report"})
	}

	logReportActivity("Schedule", scheduled.Name+" scheduled", "Auto", currentAdminName(c))

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"scheduled": scheduled})
}

// UpdateScheduledReport edits a scheduled report, including enabling or
// disabling it. Only the fields present in the request are changed.
func UpdateScheduledReport(c *fiber.Ctx) error {
	var scheduled models.ScheduledReport
	if err := config.DB.Where("id = ?", c.Params("id")).First(&scheduled).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Scheduled report not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load scheduled report"})
	}

	var input models.ScheduledReportInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	if name := strings.TrimSpace(input.Name); name != "" {
		scheduled.Name = name
	}
	if reportType := strings.TrimSpace(input.Type); reportType != "" {
		if !allowedReportTypes[reportType] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unknown report type"})
		}
		scheduled.Type = reportType
	}
	if frequency := strings.TrimSpace(input.Frequency); frequency != "" {
		scheduled.Frequency = frequency
	}
	if format := strings.TrimSpace(input.Format); format != "" {
		if !allowedReportFormats[format] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format must be CSV"})
		}
		scheduled.Format = format
	}
	if input.Enabled != nil {
		scheduled.Enabled = *input.Enabled
	}

	if err := config.DB.Save(&scheduled).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update scheduled report"})
	}

	return c.JSON(fiber.Map{"scheduled": scheduled})
}

// DeleteScheduledReport removes a scheduled report.
func DeleteScheduledReport(c *fiber.Ctx) error {
	var scheduled models.ScheduledReport
	if err := config.DB.Where("id = ?", c.Params("id")).First(&scheduled).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Scheduled report not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load scheduled report"})
	}

	if err := config.DB.Delete(&scheduled).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete scheduled report"})
	}

	return c.JSON(fiber.Map{"message": "Scheduled report deleted successfully"})
}

// ---------------------------------------------------------------------------
// Export center
// ---------------------------------------------------------------------------

// GetExportCounts returns the live record count behind each export option.
func GetExportCounts(c *fiber.Ctx) error {
	var students, activities, credits, certificates int64

	// The credits export is one row per student (a per-student credit summary),
	// so its count is the student total — not the number of approved certificates.
	queries := []struct {
		model any
		where string
		into  *int64
	}{
		{model: &models.Student{}, into: &students},
		{model: &models.Activity{}, into: &activities},
		{model: &models.Student{}, into: &credits},
		{model: &models.Certificate{}, into: &certificates},
	}

	for _, q := range queries {
		query := config.DB.Model(q.model)
		if q.where != "" {
			query = query.Where(q.where)
		}
		if err := query.Count(q.into).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load export counts"})
		}
	}

	return c.JSON(fiber.Map{
		"students":     students,
		"activities":   activities,
		"credits":      credits,
		"certificates": certificates,
	})
}

// ExportData streams a direct CSV dump of a platform dataset.
func ExportData(c *fiber.Ctx) error {
	exportType := strings.TrimSpace(c.Query("type"))
	if !allowedExportTypes[exportType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type must be one of students, activities, credits or certificates",
		})
	}

	rows, err := buildExportData(exportType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to export data"})
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.WriteAll(rows); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to export data"})
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to export data"})
	}

	logReportActivity("Export", "Exported "+exportType+" data", "CSV", currentAdminName(c))

	fileName := fmt.Sprintf("%s_export_%d.csv", exportType, time.Now().UnixNano())
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	return c.Send(buf.Bytes())
}

// buildExportData produces the full CSV dump for an export type.
func buildExportData(exportType string) ([][]string, error) {
	switch exportType {
	case "students":
		aggregates, err := buildStudentAggregates(models.GenerateReportInput{})
		if err != nil {
			return nil, err
		}
		rows := [][]string{{"Roll No", "Name", "Course", "Semester", "Credits Earned", "Activities", "Pending Certificates"}}
		for _, a := range aggregates {
			rows = append(rows, []string{a.RollNo, a.Name, a.CourseName, strconv.Itoa(a.Semester),
				strconv.Itoa(a.Credits), strconv.Itoa(a.ActivityCount), strconv.Itoa(a.Pending)})
		}
		return rows, nil

	case "credits":
		aggregates, err := buildStudentAggregates(models.GenerateReportInput{})
		if err != nil {
			return nil, err
		}
		rows := [][]string{{"Roll No", "Name", "Course", "Semester", "Credits Earned"}}
		for _, a := range aggregates {
			rows = append(rows, []string{a.RollNo, a.Name, a.CourseName, strconv.Itoa(a.Semester), strconv.Itoa(a.Credits)})
		}
		return rows, nil

	case "activities":
		var activities []models.Activity
		if err := config.DB.Order("activity_date desc").Find(&activities).Error; err != nil {
			return nil, err
		}
		rows := [][]string{{"Name", "Category", "Credits", "Mode", "Status", "Activity Date"}}
		for _, a := range activities {
			rows = append(rows, []string{a.Name, a.Category, strconv.Itoa(a.Credits), a.Mode, a.Status,
				a.ActivityDate.Format("2006-01-02")})
		}
		return rows, nil

	case "certificates":
		var certificates []models.Certificate
		if err := config.DB.Order("created_at desc").Find(&certificates).Error; err != nil {
			return nil, err
		}
		rows := [][]string{{"Student", "Activity", "Category", "Status", "Credits", "Activity Date"}}
		for _, cert := range certificates {
			rows = append(rows, []string{cert.StudentRollNo, cert.ActivityName, cert.ActivityCategory, cert.Status,
				strconv.Itoa(cert.Credits), cert.ActivityDate.Format("2006-01-02")})
		}
		return rows, nil

	default:
		return nil, fmt.Errorf("unsupported export type %q", exportType)
	}
}

// ---------------------------------------------------------------------------
// Audit log
// ---------------------------------------------------------------------------

// GetReportAuditLog returns report activity entries newest first, with
// pagination and an optional category filter.
func GetReportAuditLog(c *fiber.Ctx) error {
	limit, offset := paginate(c)

	query := config.DB.Model(&models.ReportAuditLog{})
	if category := strings.TrimSpace(c.Query("category")); category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load audit log"})
	}

	var entries []models.ReportAuditLog
	if err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&entries).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load audit log"})
	}

	return c.JSON(fiber.Map{"logs": entries, "total": total, "limit": limit, "offset": offset})
}

// ---------------------------------------------------------------------------
// Institutional overview
// ---------------------------------------------------------------------------

// GetInstitutionalOverview returns the headline figures, per-department credit
// averages and a six-month participation trend for the Reports Center's
// Institutional Overview panel. Everything is computed from live platform data.
func GetInstitutionalOverview(c *fiber.Ctx) error {
	var registered, active, activities, totalCredits int64

	if err := config.DB.Model(&models.Student{}).Count(&registered).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load overview"})
	}
	if err := config.DB.Model(&models.Student{}).Where("is_verified = ?", true).Count(&active).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load overview"})
	}
	if err := config.DB.Model(&models.Activity{}).Count(&activities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load overview"})
	}
	if err := config.DB.Model(&models.Certificate{}).Where("status = ?", "Approved").
		Select("COALESCE(SUM(credits), 0)").Scan(&totalCredits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load overview"})
	}

	avgCredits := 0.0
	if registered > 0 {
		avgCredits = float64(totalCredits) / float64(registered)
	}

	departmentScores, err := departmentScores()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load overview"})
	}

	months, activityTrend, certificateTrend, err := participationTrend()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load overview"})
	}

	return c.JSON(fiber.Map{
		"stats": fiber.Map{
			"registered_students":  registered,
			"active_students":      active,
			"activities_conducted": activities,
			"avg_credits":          avgCredits,
		},
		"department_scores": departmentScores,
		"trend": fiber.Map{
			"months":       months,
			"activities":   activityTrend,
			"certificates": certificateTrend,
		},
	})
}

// departmentScores returns the average earned credits per student for every
// canonical programme, in the canonical order. Programmes with no students are
// included with a score of 0 so the chart always shows the full course list.
func departmentScores() ([]fiber.Map, error) {
	aggregates, err := buildStudentAggregates(models.GenerateReportInput{})
	if err != nil {
		return nil, err
	}

	type dept struct {
		students int
		credits  int
	}
	depts := map[string]*dept{}
	for _, a := range aggregates {
		d, ok := depts[a.CourseName]
		if !ok {
			d = &dept{}
			depts[a.CourseName] = d
		}
		d.students++
		d.credits += a.Credits
	}

	scores := make([]fiber.Map, 0, len(models.CanonicalCourses))
	for _, name := range models.CanonicalCourses {
		score := 0
		if d, ok := depts[name]; ok && d.students > 0 {
			score = int(math.Round(float64(d.credits) / float64(d.students)))
		}
		scores = append(scores, fiber.Map{"dept": name, "score": score})
	}
	return scores, nil
}

// participationTrend returns the last six months (oldest first) with the number
// of activities and certificates created in each, keyed by created_at.
func participationTrend() (months []string, activityCounts, certificateCounts []int64, err error) {
	now := time.Now().UTC()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	for i := 5; i >= 0; i-- {
		monthStart := currentMonthStart.AddDate(0, -i, 0)
		monthEnd := monthStart.AddDate(0, 1, 0)
		months = append(months, monthStart.Format("Jan"))

		var activityCount, certificateCount int64
		if err = config.DB.Model(&models.Activity{}).
			Where("created_at >= ? AND created_at < ?", monthStart, monthEnd).
			Count(&activityCount).Error; err != nil {
			return nil, nil, nil, err
		}
		if err = config.DB.Model(&models.Certificate{}).
			Where("created_at >= ? AND created_at < ?", monthStart, monthEnd).
			Count(&certificateCount).Error; err != nil {
			return nil, nil, nil, err
		}

		activityCounts = append(activityCounts, activityCount)
		certificateCounts = append(certificateCounts, certificateCount)
	}

	return months, activityCounts, certificateCounts, nil
}
