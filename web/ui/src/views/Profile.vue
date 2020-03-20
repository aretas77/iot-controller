<template>
  <div class="settings-page">
    <b-container>
      <b-row>
        <b-col>
          <h1 class="text-xs-center mt-4">Your Profile</h1>
          <form @submit.prevent="updateProfile()">
            <fieldset>
              <fieldset class="form-group">
                <input
                  class="form-control"
                  label="Your Name:"
                  type="text"
                  v-model="currentUser.name"
                  placeholder="Enter name"
                />
              </fieldset>

              <fieldset class="form-group">
                <input
                  class="form-control form-control-lg"
                  type="password"
                  v-model="currentUser.password"
                  placeholder="Password"
                />
              </fieldset>

              <button class="btn btn-lg btn-primary pull-xs-right">
                Update Profile
              </button>
           </fieldset>
          </form>
          <!-- Line break for logout button -->
          <hr />
          <button @click="logout" class="btn btn-outline-danger">
            Logout
          </button>
        </b-col>
        <b-col>
          <h1 class="text-xs-center mt-4">Network settings</h1>
          <form @submit.prevent="selectNetwork()" class="float-left">
            <fieldset class="form-group">
              <select v-model="network">
                <option disabled value="">Select a network</option>
              </select>
            </fieldset>
            <button class="btn btn-lg btn-primary pull-xs-right">
              Set network
            </button>
          </form>

        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import {
  UPDATE_USER,
  LOGOUT,
  FETCH_NETWORKS,
  SET_NETWORK
} from '@/store/actions.type'

// XXX: Export Profile and Network settings into modules.
export default {
  name: 'Profile',
  data () {
    return {
      network: ''
    }
  },
  computed: {
    ...mapGetters(['currentUser', 'currentNetwork'])
  },
  mounted () {
    this.fetchNetworks()
  },
  methods: {
    updateProfile () {
      this.$store.dispatch(UPDATE_USER, this.currentUser).then(() => {
        this.$router.push({ name: '/' })
      })
    },
    selectNetwork () {
      this.$store.dispatch(SET_NETWORK, null)
    },
    async logout () {
      this.$store.dispatch(LOGOUT).then(() => {
        this.$router.push({ name: '/' })
      })
    },
    fetchNetworks () {
      this.$store.dispatch(FETCH_NETWORKS, this.currentUser.id)
    }
  }
}
</script>
