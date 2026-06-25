<script lang="ts">
  import HeaderSection from './HeaderSection.svelte';
  import StudentRegistrationSection from './StudentRegistrationSection.svelte';
  import FooterSection from './FooterSection.svelte';
</script>

<svelte:head>
  <title>Student Registration | iSPARC IIPS DAVV</title>
  <meta name="description" content="Register your student account with iSPARC - International Institute of Professional Studies, DAVV, Indore." />
</svelte:head>

<div class="min-h-screen bg-bg-main flex flex-col items-center w-full">
  <HeaderSection />
  <main class="flex-grow w-full flex justify-center items-start">
    <StudentRegistrationSection />
  </main>
  <FooterSection />
</div>
