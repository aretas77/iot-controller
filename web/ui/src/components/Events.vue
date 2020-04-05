<template>
  <b-table striped hover :items="events"></b-table>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'NodeEvents',
  props: {
    events: {
      type: Array,
      required: true
    }
  },

  computed: {
    ...mapGetters(['isAuthenticated'])
  }
}
</script>
