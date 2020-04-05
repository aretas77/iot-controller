<template>
  <div>
    <b-table
      striped hover
      id="eventTable"
      :items="events"
      :busy="busy"
      :show-empty="true"
      :per-page="perPage"
      :current-page="currentPage"
    >
      <template v-slot:table-busy>
        <div class="text-center text-danger my-2">
          <b-spinner class="align-middle"></b-spinner>
          <strong>Loading...</strong>
        </div>
      </template>
    </b-table>
    <b-row class="m-2" align-h="center">
      <b-pagination
        v-model="currentPage"
        :total-rows="rows"
        :per-page="perPage"
        aria-controls="eventTable"
        >
      </b-pagination>
    </b-row>
  </div>
</template>

<script>
export default {
  name: 'NodeEvents',
  props: {
    events: {
      type: Array,
      required: true
    },
    eventsLen: {
      type: Number,
      required: true
    },
    busy: {
      type: Boolean,
      required: true
    },
    itemsPerPage: {
      type: Number,
      required: false,
      default: 10
    }
  },
  data () {
    return {
      currentPage: 1,
      perPage: this.itemsPerPage
    }
  },
  computed: {
    rows () {
      return this.eventsLen
    }
  }
}
</script>
