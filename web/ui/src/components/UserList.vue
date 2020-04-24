<template>
  <b-container>
    <b-row class="m-3">
      <b-col>
        <b-table
          striped hover
          id="userTable"
          :items="users"
          :show-empty="true"
          :busy="busy"
          >
          <template v-slot:table-busy>
            <div class="text-center text-danger my-2">
              <b-spinner class="align-middle"></b-spinner>
              <strong>Loading...</strong>
            </div>
          </template>
        </b-table>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
export default {
  name: 'UserList',
  props: {
    users: {
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
