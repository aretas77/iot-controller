<template>
  <b-table striped hover :items="events" :busy="busy" :show-empty="true">
    <template v-slot:table-busy>
      <div class="text-center text-danger my-2">
        <b-spinner class="align-middle"></b-spinner>
        <strong>Loading...</strong>
      </div>
    </template>
  </b-table>
</template>

<script>
export default {
  name: 'NodeEvents',
  props: {
    events: {
      type: Array,
      required: true
    },
    busy: {
      type: Boolean,
      required: true
    }
  }
}
</script>
