<template>
  <div class="settings-page">
    <b-container>
      <b-row>
        <b-col>
          <h1 class="text-xs-center mt-4">Your Profile</h1>
          <b-form @submit.prevent="updateProfile()">
            <b-form-group>
              <b-form-input
                class="form-control"
                label="Your Name:"
                type="text"
                v-model="currentUser.name"
                placeholder="Enter name"
                />
            </b-form-group>

            <b-form-group>
              <b-form-input
                class="form-control form-control-lg"
                type="password"
                v-model="password"
                placeholder="Current password"
                />
            </b-form-group>
            <b-form-group>
              <b-form-input
                class="form-control form-control-lg"
                type="password"
                v-model="password"
                placeholder="Password"
                />
            </b-form-group>

            <b-form-group>
              <button class="btn btn-lg btn-primary pull-xs-right">
                Update Profile
              </button>
            </b-form-group>
          </b-form>
          <!-- Line break for logout button -->
          <hr />
          <button @click="logout" class="btn btn-outline-danger">
            Logout
          </button>
        </b-col>
        <b-col>
          <h1 class="text-xs-center mt-4">Network settings</h1>
          <b-form @submit.prevent="selectNetwork()">
            <b-form-group>
              <b-form-select
                v-model="network"
                :options="networksOptions"
                >
                <template v-slot:first v-if="currentNetwork.name">
                  <b-form-select-option
                    :value="null"
                    disabled
                    >
                    -- Current network: {{ currentNetwork.name }} --
                  </b-form-select-option>
                </template>
              </b-form-select>
            </b-form-group>
            <button class="btn btn-lg btn-primary pull-xs-right">
              Set network
            </button>
          </b-form>

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
  UPDATE_CURRENT_NETWORK
} from '@/store/actions.type'

// XXX: Export Profile and Network settings into modules.
export default {
  name: 'Profile',
  data () {
    return {
      network: null,
      password: '',
      networksOptions: []
    }
  },
  computed: {
    ...mapGetters(['currentUser', 'currentNetwork', 'networks'])
  },
  watch: {
    networks (newValue, oldValue) {
      for (const key of newValue) {
        this.networksOptions.push({ text: key.name, value: key })
      }
    }
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
    // will set the current global network as the selected network.
    selectNetwork () {
      this.$store.dispatch(UPDATE_CURRENT_NETWORK, this.network)
    },
    async logout () {
      this.$store.dispatch(LOGOUT).then(() => {
        this.$router.push({ name: '/' })
      })
    },
    fetchNetworks () {
      this.$store.dispatch(FETCH_NETWORKS, this.currentUser.ID)
    }
  }
}
</script>
