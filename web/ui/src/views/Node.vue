<template>
  <div class="node-page mt-3">
    <div class="banner">
      <div class="banner-info">
        <div class="py-4">
          <h1>{{ node.name }}</h1>
          <p class="mb-0">{{ node.mac }}</p>
        </div>
      </div>
    </div>
    <div>
      <b-tabs content-class="mt-3" class="m-2" @input="onChangedTab">

        <b-tab title="Statistics" active lazy>
          <b-container v-if="isLoadingStats">
            Loading...
          </b-container>

          <b-container fluid class="w-100 p-3" v-else>
            <b-row align-h="between" class="m-4">
              <b-col cols="5">
                <TemperatureChart :entries='statsEntries'></TemperatureChart>
              </b-col>
              <b-col cols="5">
                <SensorReadingsFreq />
              </b-col>
            </b-row>
          </b-container>

        </b-tab>

        <b-tab title="Models" lazy>
          <b-container v-if="isLoadingEvents" fluid class="w-100 p-3">
            <NodeEvents :events='events' :busy='isLoadingEvents' :eventsLen='eventsCount'>
            </NodeEvents>
          </b-container>
          <b-container fluid class="w-100 p-3" v-else>
            <NodeEvents :events='events' :busy='isLoadingEvents' :eventsLen='eventsCount'>
            </NodeEvents>
          </b-container>
        </b-tab>

        <b-tab title="Information" lazy>
          <b-container fluid class="w-100 p-3">
            <b-row align-h="between" class="m-4">
              <b-col cols>
                <BatteryLevelsChart
                  :interval='60'
                  :entries='parsedBatteryLevels'
                >
                </BatteryLevelsChart>
              </b-col>
            </b-row>
          </b-container>
        </b-tab>
      </b-tabs>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import store from '@/store'
import marked from 'marked'
import TemperatureChart from '@/components/TemperatureChart'
import BatteryLevelsChart from '@/components/BatteryLevelsChart'
import SensorReadingsFreq from '@/components/SensorReadingsFreq'
import NodeEvents from '@/components/Events'
import {
  FETCH_NODE,
  FETCH_NODE_STATS,
  FETCH_NODE_EVENTS,
  CHECK_AUTH
} from '@/store/actions.type'

export default {
  name: 'iotctl-node',
  components: {
    TemperatureChart,
    BatteryLevelsChart,
    SensorReadingsFreq,
    NodeEvents
  },
  props: {
    slug: {
      type: [Number, String],
      required: true
    }
  },
  data () {
    return {
      activeTab: 0,
      parsedBatteryLevels: []
    }
  },
  // actions
  beforeRouteEnter (to, from, next) {
    Promise.all([
      store.dispatch(FETCH_NODE, to.params.slug)
    ]).then(() => {
      next()
    })
  },
  computed: {
    ...mapGetters(['node', 'currentUser', 'isAuthenticated', 'isLoadingStats',
      'statsEntries', 'events', 'isLoadingEvents', 'eventsCount'])
  },
  methods: {
    parseMarkdown (content) {
      return marked(content)
    },
    fetchNodeStatistics () {
      this.$store.dispatch(FETCH_NODE_STATS, this.node.ID)
    },
    fetchNodeEvents () {
      this.$store.dispatch(FETCH_NODE_EVENTS, this.node.ID)
    },
    parseBatteryLevels () {
      var tempParsedBatteryLevels = this.statsEntries.map(entry => {
        var beforeHours = this.$moment().subtract(12, 'h')
        var afterHours = this.$moment().add(12, 'h')

        // Check if current 24hrs
        if (this.$moment(entry.temp_read_time).isBetween(beforeHours, afterHours)) {
          return {
            x: this.$moment(entry.temp_read_time, this.$moment.ISO_8601).format('hh:00 a'),
            y: entry.battery_left_per + 10
          }
        }
      })

      tempParsedBatteryLevels = tempParsedBatteryLevels.filter(function (entry) {
        return entry !== undefined
      })

      /* eslint-disable no-unused-vars */
      const newLevels = tempParsedBatteryLevels.reduce(function (result, item) {
        if (item.x in result) {
          result[item.x].total += item.y
          result[item.x].count++
        } else {
          result[item.x] = {
            total: item.y,
            count: 1
          }
        }
        return result
      }, {})

      this.parsedBatteryLevels = []
      for (const [key, value] of Object.entries(newLevels)) {
        this.parsedBatteryLevels.push({ x: key, y: value.total / value.count })
      }

      /* eslint-enable no-unused-vars */
    },
    onChangedTab (tabIndex) {
      this.$store.dispatch(CHECK_AUTH)

      if (tabIndex === 0) {
        // refresh statistics
        // TODO: could probably just check with the server if there are new
        // updates
        this.fetchNodeStatistics()
      } else if (tabIndex === 1) {
        this.fetchNodeEvents()
      } else if (tabIndex) {
        this.fetchNodeStatistics()
        this.parseBatteryLevels()
      }
    }
  },
  mounted () {
    this.fetchNodeStatistics()
    this.fetchNodeEvents()
    this.parseBatteryLevels()
  }
}
</script>

<style lang="scss">
.banner {
  background: #5cb85c;
  box-shadow: inset 0 8px 8px -8px rgba(0,0,0,.3),inset 0 -8px 8px -8px rgba(0,0,0,.3);
  color: #fff;
  margin-bottom: 1rem;
  margin-top: 1rem;
}
.banner-info {
  text-shadow: 0 1px 3px rgba(0,0,0,.3);
  font-weight: 700 !important;

}
</style>
