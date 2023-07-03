<template>
<div class="modes">
  <div class="mode" v-for="mode in connection.data.modes">
    {{ mode }}
  </div>
</div>
</template>

<script src="settings.js"></script>

<style scoped lang="sass" src="settings.sass"></style>