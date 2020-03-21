<template>
  <div class="node-editor">
    <b-container class="border mt-3">
      <b-row>
        <b-col>
          <IotctlListErrors :errors="errors" />
          <b-form
            id="register_node"
            @submit.prevent="onPublish(unregisteredNode.slug)"
            >

            <fieldset :disabled="inProgress">
              <b-form-group
                id="mac-group"
                label="Node MAC address"
                label-for="mac-input"
                >
                <b-form-input
                  id="mac-input"
                  type="text"
                  class="form-control form-control-lg"
                  v-model="unregisteredNode.mac"
                  :state="validateMAC"
                />
                <b-form-invalid-feedback :state="validateMAC">
                  Invalid MAC. The format is AA:BB:CC:DD:EE:FF.
                </b-form-invalid-feedback>
                <b-form-valid-feedback :state="validateMAC">
                  Looks good.
                </b-form-valid-feedback>
              </b-form-group>

              <b-form-group>
                <b-form-select
                  v-model="network"
                  :options="networksOptions"
                  :state="validateNetwork"
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
                <b-form-invalid-feedback :state="validateNetwork">
                  Please select a network.
                </b-form-invalid-feedback>
              </b-form-group>
            </fieldset>

            <!-- Form submit button -->
            <button
              :disabled="!validateMAC || !validateNetwork"
              class="btn btn-lg pull-xs-right btn-primary m-3"
              type="submit"
            >
              Register Node
            </button>
          </b-form>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import store from '@/store'
import IotctlListErrors from '@/components/ListErrors'
import {
  UNREGISTERED_NODE_ADD,
  NODE_RESET_STATE,
  FETCH_NETWORKS
} from '@/store/actions.type'

export default {
  name: 'NodeCreate',
  components: {
    IotctlListErrors
  },
  props: {
    previousNode: {
      type: Object,
      required: false
    }
  },
  mounted () {
    this.fetchNetworks()
  },
  watch: {
    // If there was a change in `networks` we build a new option list.
    networks (newValue, oldValue) {
      for (const key of newValue) {
        this.networksOptions.push({ text: key.name, value: key })
      }
    }
  },
  async beforeRouterUpdate (to, from, next) {
    // Reset state if user goes from /editor/:id to /editor
    // The component is not recreated so we use hook to reset the state.
    await store.dispatch(NODE_RESET_STATE)
    return next()
  },

  async beforeRouteEnter (to, from, next) {
    await store.dispatch(NODE_RESET_STATE)
    if (to.params.slug !== undefined) {
    }
    return next()
  },

  data () {
    return {
      inProgress: false,
      errors: {},
      unregisteredNode: {},
      network: null,
      networksOptions: [],
      validated: false
    }
  },

  computed: {
    ...mapGetters(['currentUser', 'currentNetwork', 'networks']),
    validateMAC () {
      return this.validMAC(this.unregisteredNode.mac)
    },
    validateNetwork () {
      return this.validNetwork()
    }
  },
  methods: {
    // Publish the filled up form of Node creation
    onPublish (slug) {
      const action = UNREGISTERED_NODE_ADD
      this.inProgress = true
      console.log(this.network)
      this.$store
        .dispatch(action, {
          mac: this.unregisteredNode.mac,
          network_refer: this.network.ID
        }).then(({ data }) => {
          this.inProgress = false
          this.$router.push({
            name: 'node',
            params: { slug: data.unregisteredNode.slug }
          })
        }).catch(({ response }) => {
          this.inProgress = false
          // this.errors = response.data.errors
        })
    },
    fetchNetworks () {
      this.$store.dispatch(FETCH_NETWORKS, this.currentUser.ID)
    },
    validMAC: function (mac) {
      var re = /^((([a-fA-F0-9][a-fA-F0-9]+[-]){5}|([a-fA-F0-9][a-fA-F0-9]+[:]){5})([a-fA-F0-9][a-fA-F0-9])$)|(^([a-fA-F0-9][a-fA-F0-9][a-fA-F0-9][a-fA-F0-9]+[.]){2}([a-fA-F0-9][a-fA-F0-9][a-fA-F0-9][a-fA-F0-9]))$/
      return re.test(mac)
    },
    validNetwork: function () {
      if (this.network === null || this.network === undefined) {
        return false
      }
      for (const key of this.networks) {
        if (this.currentNetwork.name === key.name) {
          return true
        }
      }
      return false
    }
  }
}
</script>

<style lang="scss">
.register_node {
  position: relative;
}
</style>
