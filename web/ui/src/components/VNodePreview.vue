<template>
  <b-container fluid class="node-preview-container">
    <b-row fluid class="m-3 node-preview-row">

      <b-col cols="2" class="d-flex align-items-start">
        <router-link :to="nodeLink" class="preview-link">
          <h3 v-text="node.name" />
          <p v-text="node.mac" />
        </router-link>
      </b-col>

      <b-col cols="3" class="text-left">
        <span class="date">Last active: {{ node.last_received | time }}</span>
        <br />
        <span>Location: {{ node.location }}</span>
        <br />
        <span>IP Address: {{ node.ipv4 }}</span>
        <br />
        <span class="battery" v-bind:class="activeColor">
          Battery: {{ node.battery_left_per | percentage }}% ({{ node.battery_left_mah }} / {{ node.battery_total_mah }} mAh)
        </span>
      </b-col>

      <b-col offset="4" class="">
        <NodeMeta isPreview :node="node" :actions="true" />
      </b-col>

    </b-row>
    <hr/>
  </b-container>
</template>

<script>
import NodeMeta from './NodeMeta'

export default {
  name: 'IotctlNodePreview',
  components: {
    NodeMeta
  },
  data () {
    return {
      good: true,
      average: false,
      bad: false
    }
  },
  props: {
    node: { type: Object, required: true }
  },
  computed: {
    nodeLink () {
      return {
        name: 'node',
        params: {
          slug: this.node.ID
        }
      }
    },
    activeColor: function () {
      return {
        'battery-good': this.good && !this.average && !this.bad,
        'battery-average': !this.good && this.average && !this.bad,
        'battery-bad': !this.good && !this.average && this.bad
      }
    }
  },
  watch: {
    node () {
      if (this.node.battery_left_per >= 70) {
        this.good = true
        this.bad = false
        this.average = false
      } else if (this.node.battery_left_per < 70 && this.node.battery_left_per > 40) {
        this.good = false
        this.average = true
        this.bad = false
      } else {
        this.good = false
        this.average = false
        this.bad = true
      }
    }
  }
}
</script>

<style lang="scss">
.node-preview-container {
  margin: 7px;
}

.node-preview-row {
}

.preview-link {
}

.battery-bad {
  color: red;
}
.battery-good {
  color: green;
}
.batter-average {
  color: yellow;
}
</style>
