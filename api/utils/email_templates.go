package utils

import (
	"fmt"
	"html"
	"strings"
)

// Email templates.
//
// Written as table-based HTML with inline styles on purpose: email clients
// (Gmail, Outlook, Apple Mail) strip <style> blocks and ignore flexbox/grid, so
// the layout patterns that look outdated for the web are the ones that render
// reliably in an inbox. Width is capped at 600px, the standard safe width.
//
// Every message is sent as multipart: a plain-text part for clients with HTML
// disabled, and this HTML part for everyone else.

// Brand colours, matching the portal.
const (
	brandNavy   = "#0B1535"
	brandMaroon = "#881B1B"
	brandBorder = "#D8DDE5"
	brandMuted  = "#6B7280"
	brandBg     = "#F7F6F3"
)

// emailShell wraps body content in the shared header/footer chrome.
//
// preheader is the short summary line inboxes show next to the subject; it is
// hidden in the message itself.
func emailShell(preheader, heading, bodyHTML string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%[2]s</title>
</head>
<body style="margin:0; padding:0; background-color:%[6]s; font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif; -webkit-font-smoothing:antialiased;">
  <div style="display:none; font-size:1px; color:%[6]s; line-height:1px; max-height:0; max-width:0; opacity:0; overflow:hidden;">%[1]s</div>
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:%[6]s; padding:24px 12px;">
    <tr>
      <td align="center">
        <table role="presentation" width="600" cellpadding="0" cellspacing="0" border="0" style="width:600px; max-width:600px; background-color:#ffffff; border:1px solid %[5]s; border-radius:12px; overflow:hidden;">

          <!-- Header -->
          <tr>
            <td style="background-color:%[3]s; padding:28px 32px;">
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0">
                <tr>
                  <td>
                    <div style="font-size:22px; font-weight:700; color:#ffffff; letter-spacing:0.5px;">iSPARC</div>
                    <div style="font-size:11px; color:#B9C1D4; letter-spacing:1.5px; text-transform:uppercase; margin-top:4px;">IIPS Skill, Personality Advancement &amp; Refinement Cell</div>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          <!-- Body -->
          <tr>
            <td style="padding:32px;">
              <h1 style="margin:0 0 16px 0; font-size:20px; line-height:1.35; color:%[3]s; font-weight:700;">%[2]s</h1>
              %[4]s
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="border-top:1px solid %[5]s; padding:20px 32px 24px 32px;">
              <p style="margin:0 0 6px 0; font-size:12px; color:%[7]s; line-height:1.6;">
                This is an automated message from the iSPARC portal. Please do not reply to it.
              </p>
              <p style="margin:0; font-size:12px; color:%[7]s; line-height:1.6;">
                International Institute of Professional Studies, DAVV, Indore
              </p>
            </td>
          </tr>

        </table>

        <p style="margin:16px 0 0 0; font-size:11px; color:%[7]s;">
          You received this email because an action was requested for your iSPARC account.
        </p>
      </td>
    </tr>
  </table>
</body>
</html>`, html.EscapeString(preheader), html.EscapeString(heading), brandNavy, bodyHTML, brandBorder, brandBg, brandMuted)
}

// otpEmail renders a verification-code message. intro explains why the code was
// sent and action names what it unlocks, so the same layout serves both account
// verification and password reset.
func otpEmail(heading, intro, code, action string) string {
	body := fmt.Sprintf(`
              <p style="margin:0 0 20px 0; font-size:15px; color:#374151; line-height:1.65;">%s</p>

              <!-- Code block -->
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 20px 0;">
                <tr>
                  <td align="center" style="background-color:%s; border:1px solid %s; border-radius:10px; padding:24px 16px;">
                    <div style="font-size:11px; font-weight:700; color:%s; letter-spacing:1.5px; text-transform:uppercase; margin-bottom:10px;">Your verification code</div>
                    <div style="font-size:34px; font-weight:700; color:%s; letter-spacing:10px; font-family:'SFMono-Regular',Consolas,'Liberation Mono',Menlo,monospace; padding-left:10px;">%s</div>
                  </td>
                </tr>
              </table>

              <p style="margin:0 0 20px 0; font-size:14px; color:#374151; line-height:1.65;">
                Enter this code in the portal to %s. The code expires in <strong>15 minutes</strong> and can be used once.
              </p>

              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin:0;">
                <tr>
                  <td style="background-color:#FFFBEB; border-left:3px solid #D97706; border-radius:4px; padding:12px 14px;">
                    <p style="margin:0; font-size:13px; color:#92400E; line-height:1.6;">
                      If you did not request this, you can ignore this email &mdash; nothing will change on your account. Never share this code with anyone.
                    </p>
                  </td>
                </tr>
              </table>`,
		html.EscapeString(intro), brandBg, brandBorder, brandMuted, brandMaroon,
		html.EscapeString(code), html.EscapeString(action))

	return emailShell(fmt.Sprintf("Your iSPARC verification code is %s", code), heading, body)
}

// noticeEmail renders a general message (mentor notices, activity reminders).
// The body is plain text from the caller; newlines become paragraphs.
func noticeEmail(heading, message string) string {
	var paragraphs strings.Builder
	for _, line := range strings.Split(message, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		fmt.Fprintf(&paragraphs,
			`<p style="margin:0 0 14px 0; font-size:15px; color:#374151; line-height:1.65;">%s</p>`,
			html.EscapeString(trimmed),
		)
	}

	preheader := message
	if len(preheader) > 120 {
		preheader = preheader[:120]
	}
	return emailShell(preheader, heading, paragraphs.String())
}
