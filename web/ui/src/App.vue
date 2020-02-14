<template>
  <div id="app">
    <b-navbar toggleable="md" type="dark" variant="dark">
        <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
        <b-navbar-brand to="/">IoT Controller</b-navbar-brand>
        <b-collapse is-nav id="nav_collapse">
            <b-navbar-nav>
                <b-nav-item to="/">Home</b-nav-item>
                <b-nav-item to="/nodes">Nodes</b-nav-item>
                <b-nav-item href="#" @click.prevent="login" v-if="!activeUser">Login</b-nav-item>
                <b-nav-item href="#" @click.prevent="logout" v-else>Logout</b-nav-item>
            </b-navbar-nav>
        </b-collapse>
    </b-navbar>
    <div id="nav">
      <router-link to="/nodes">Nodes</router-link> |
      <router-link to="/models">Models</router-link>
    </div>
    <router-view/>
  </div>
</template>

<script>

export default {
    name: 'app',
    data () {
        return {
            activeUser: null
        }
    },
    async created() {
        await this.refreshActiveUser()
    },
    watch: {
        // everytime a route is changed refresh the activeUser
        '$route': 'refreshActiveUser'
    },
    methods: {
        login() {
            this.$auth.loginRedirect()
        },
    }
}

</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
