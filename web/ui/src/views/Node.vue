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
                <TemperatureChart :entries='parsedTemperature'></TemperatureChart>
              </b-col>
              <b-col cols="5">
                <SensorReadingsFreq
                  :entries='parsedSendFrequency'
                  :entriesFrames='parsedSendFrequencyFrames'
                >
                </SensorReadingsFreq>
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

        <b-tab title="Battery status" lazy>
          <b-container fluid class="w-100 p-3">
            <b-row>
              <b-col cols="3" class="text-left">
                <span>
                  Battery: {{ node.battery_left_per | percentage }}%
                </span>
                <br />
                <span>
                  Battery (mAh): {{ node.battery_left_mah }} / {{ node.battery_total_mah}} mAh
                </span>
              </b-col>
            </b-row>

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
      parsedTemperature: [],
      parsedBatteryLevels: [],
      parsedSendFrequency: [],
      parsedSendFrequencyFrames: []
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
    /**
      *countValuesAtHour is used to calculate total values in a given hour.
      * Given:
      * 0: Object { x: "02:00 am", y: 13 }
      * 1: Object { x: "02:00 am", y: 14 }
      * 2: Object { x: "02:00 am", y: 15 }
      * 3: Object { x: "02:00 am", y: 16 }
      * 4: Object { x: "02:00 am", y: 17 }
      * 5: Object { x: "02:00 am", y: 18 }
      * 6: Object { x: "03:00 am", y: 19 }
      * 7: Object { x: "03:00 am", y: 20 }
      * Result:
      * "02:00 am": Object { total: 93, count: 6 }
      * "03:00 am": Object { total: 39, count: 2 }
     */
    countValuesAtHour (obj) {
      const newLevels = obj.reduce(function (result, item) {
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
      return newLevels
    },
    fetchNodeStatistics () {
      this.$store.dispatch(FETCH_NODE_STATS, this.node.ID)
    },
    fetchNodeEvents () {
      this.$store.dispatch(FETCH_NODE_EVENTS, this.node.ID)
    },
    parseSendFrequency () {
      var tmpParsedSendFrequency = this.statsEntries.map(entry => {
        var beforeHours = this.$moment().subtract(12, 'h')
        var afterHours = this.$moment().add(12, 'h')
        var currentTime = this.$moment(entry.temp_read_time)

        // Check if current 24hrs
        if (currentTime.isBetween(beforeHours, afterHours)) {
          return {
            x: this.$moment(entry.temp_read_time, this.$moment.ISO_8601).format('hh:00 a'),
            y: entry.send_times
          }
        }
      })

      tmpParsedSendFrequency = tmpParsedSendFrequency.filter(function (entry) {
        return entry !== undefined
      })

      var newLevels = this.countValuesAtHour(tmpParsedSendFrequency)
      console.log(newLevels)
      var prevValue = 0

      this.parsedSendFrequency = []
      this.parsedSendFrequencyFrames = []
      for (const [key, value] of Object.entries(newLevels)) {
        // accumulate send counters
        this.parsedSendFrequency.push({ x: key, y: value.count + prevValue })
        this.parsedSendFrequencyFrames.push({ x: key, y: value.count })
        prevValue += value.count
      }
    },
    parseBatteryLevels () {
      var tempParsedBatteryLevels = this.statsEntries.map(entry => {
        var beforeHours = this.$moment().subtract(12, 'h')
        var afterHours = this.$moment().add(12, 'h')

        // Check if current 24hrs
        if (this.$moment(entry.temp_read_time).isBetween(beforeHours, afterHours)) {
          return {
            x: this.$moment(entry.temp_read_time, this.$moment.ISO_8601).format('hh:00 a'),
            y: entry.battery_left_per
          }
        }
      })

      tempParsedBatteryLevels = tempParsedBatteryLevels.filter(function (entry) {
        return entry !== undefined
      })

      /* eslint-disable no-unused-vars */
      var newLevels = this.countValuesAtHour(tempParsedBatteryLevels)

      this.parsedBatteryLevels = []
      for (const [key, value] of Object.entries(newLevels)) {
        this.parsedBatteryLevels.push({ x: key, y: value.total / value.count })
      }

      /* eslint-enable no-unused-vars */
    },
    parseTemperature () {
      var tmpParsedTemperature = this.statsEntries.map(entry => {
        return {
          x: this.$moment(entry.temp_read_time, this.$moment.ISO_8601).format('HH:mm'),
          y: entry.temperature
        }
      })

      tmpParsedTemperature = tmpParsedTemperature.filter(function (entry) {
        return entry !== undefined
      })

      this.parsedTemperature = tmpParsedTemperature
    },
    onChangedTab (tabIndex) {
      this.$store.dispatch(CHECK_AUTH)

      if (tabIndex === 0) {
        // refresh statistics
        // TODO: could probably just check with the server if there are new
        // updates
        this.fetchNodeStatistics()
        this.parseTemperature()
        this.parseSendFrequency()
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
  },
  watch: {
    statsEntries () {
      this.parseSendFrequency()
      this.parseTemperature()
    }
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
