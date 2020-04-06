<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <HelloWorld />
  </div>
</template>

<script>
// @ is an alias to /src
import { mapGetters } from 'vuex'
import HelloWorld from '@/components/HelloWorld.vue'

export default {
  name: 'Home',
  components: {
    HelloWorld
  },
  computed: {
    ...mapGetters(['isAuthenticated'])
  }
}
</script>
