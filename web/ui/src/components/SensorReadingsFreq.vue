<script>
import { Bar } from 'vue-chartjs'

export default {
  extends: Bar,
  props: {
    // sensor reading events cumulated from each hour before
    entries: {
      type: Array,
      required: true
    },
    // sensor reading events each hour but not cumulated
    entriesFrames: {
      type: Array,
      required: false
    }
  },

  data () {
    return {
      datacollection: {
        labels: [],
        datasets: []
      },
      options: {}
    }
  },

  methods: {
    generateLabels (interval) {
      var x = 60 // minutes interval
      var tt = 0 // start time
      var ap = ['AM', 'PM']

      for (var i = 0; tt < 24 * 60; ++i) {
        var hh = Math.floor(tt / 60)
        var mm = tt % 60
        var tmp = ('0' + (hh % 12)).slice(-2) + ':' +
          ('0' + mm).slice(-2) + ap[Math.floor(hh / 12)]
        if (tmp === '00:00AM') {
          this.datacollection.labels[i] = '12:00AM'
        } else if (tmp === '00:00PM') {
          this.datacollection.labels[i] = '12:00PM'
        } else {
          this.datacollection.labels[i] = tmp
        }
        tt = tt + x
      }
    },
    generateChart () {
      this.datacollection.datasets = [
        {
          label: 'Total send events',
          backgroundColor: '#5cb85c',
          borderWidth: 2,
          fill: true,
          type: 'bar',
          data: this.entries
        },
        {
          label: 'Send events',
          backgroundColor: 'red',
          borderWidth: 2,
          fill: true,
          type: 'bar',
          data: this.entriesFrames
        }
      ]

      this.options = {
        legend: {
          display: true
        },
        layout: {
        },
        offset: true,
        response: true,
        maintainAspectRatio: false,
        scales: {
          xAxes: [{
            type: 'time',
            time: {
              parser: 'hh:00 a',
              unitStepSize: 1,
              unit: 'hour',
              displayFormats: {
                hour: 'h A'
              },
              stacked: true
            },
            ticks: {
              source: 'labels',
              autoSkip: false
            }
          }],
          yAxes: [{
            ticks: {
              beginAtZero: true,
              min: 0,
              callback: function (value, index, values) {
                return value
              }
            }
          }]
        }
      }
    }
  },
  watch: {
    entries () {
      this.renderChart(this.datacollection, this.options)
    }
  },
  mounted () {
    // console.log(this.entriesFrames)
    this.generateLabels(60)
    this.generateChart()
    this.renderChart(this.datacollection, this.options)
  }
}
</script>
