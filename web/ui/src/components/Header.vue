<template>
  <div>
    <b-navbar class="sticky-top fixed" toggleable="md" type="dark" variant="dark">
      <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
      <b-collapse is-nav id="nav_collapse">
        <b-navbar-brand to="/">IoT Controller</b-navbar-brand>

        <b-navbar-nav>
          <b-nav-item to="/">Home</b-nav-item>
          <b-nav-item v-if="isAuthenticated" to="/nodes">Nodes</b-nav-item>
          <b-nav-item v-if="!isAuthenticated" to="/login">Login</b-nav-item>
          <b-nav-item v-if="isAuthenticated" @click="logout">Logout</b-nav-item>
        </b-navbar-nav>

        <!-- Right aligned nav items -->
        <b-navbar-nav class="ml-auto" v-if="isAuthenticated">
          <!-- Links to the currently active network -->
          <b-nav-item v-if="activeNetwork">
            <router-link
              class="nav-link"
              active-class="active"
              exact :to="{
                name: 'Home',
                params: { id: activeNetwork.id }
              }"
              >
              {{ activeNetwork.name }}
            </router-link>
          </b-nav-item>

          <!-- Links to the currently active user profile -->

          <!-- XXX: this could break if currentUser is null, but we assume
                that it is not null because we are authenticated and thus there
                should be a currentUser. -->
          <b-nav-item v-if="isAuthenticated">
            <router-link
              class="nav-link"
              active-class="active"
              exact :to="{
                name: 'profile',
                params: { username: currentUser.username }
              }"
              >
              {{ currentUser.username }}
            </router-link>
          </b-nav-item>
          <!-- Right aligned nav items end -->
        </b-navbar-nav>

      </b-collapse>
    </b-navbar>
    <router-view/>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import {
  LOGOUT,
  NODE_RESET_STATE
} from '@/store/actions.type'

export default {
  name: 'Header',
  computed: {
    ...mapGetters(['currentUser', 'isAuthenticated', 'currentNetwork'])
  },
  data () {
    return {
      activeUser: null,
      activeNetwork: null
    }
  },
  async created () {
    await this.refreshActiveNetwork()
  },
  watch: {
    // everytime a route is changed refresh the activeUser
    isAuthenticated (newValue, oldValue) {
      console.log(`isAuthenticated changed from ${oldValue} to ${newValue}`)

      if (newValue === false) {
        this.$store.dispatch(NODE_RESET_STATE)
        this.$router.push('/login')
      }
    },
    currentNetwork () {
      this.refreshActiveNetwork()
    }
  },
  methods: {
    async refreshActiveNetwork () {
      this.activeNetwork = this.currentNetwork
    },
    async logout () {
      this.$store.dispatch(LOGOUT).then(() => {
        this.$router.push('/')
      })
    }
  }
}
</script>

<style scoped lang="scss">
#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #3c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
