<script setup>
import { LxTile } from '@wntr/lx-ui';
</script>

<template>
  <div class="lx-dashboard">
    <LxTile
      icon="lock"
      label="Iet uz lapu nav tiesību"
      :to="{ name: 'forbidden' }"
    />
  </div>
</template>
