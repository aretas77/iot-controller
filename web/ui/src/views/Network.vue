<template>
  <b-container>
    <UserList />
  </b-container>
</template>

<script>
import { mapGetters } from 'vuex'
import UserList from '@/components/UserList'
import {
  FETCH_NETWORK_USERS
} from '@/store/actions.type'

export default {
  name: 'Network',
  components: {
    UserList
  },
  beforeRouteEnter (to, from, next) {
    Promise.all([
    ]).then(() => {
      next()
    })
  },
  computed: {
    ...mapGetters(['currentUser', 'currentNetwork', 'isAuthenticated'])
  },
  methods: {
    fetchNetworkUsers () {
      this.$store.dispatch(FETCH_NETWORK_USERS, this.currentNetwork.ID)
    }
  }
}
</script>
