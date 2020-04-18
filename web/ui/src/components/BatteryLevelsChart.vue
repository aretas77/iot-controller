<script>
import { Bar } from 'vue-chartjs'

export default {
  extends: Bar,

  props: {
    entries: {
      type: Array,
      required: true
    },
    interval: {
      type: Number,
      required: false,
      default: 60
    }
  },

  data () {
    return {
      datacollection: {
        labels: []
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
          label: 'Battery levels',
          backgroundColor: '#5cb85c',
          borderWidth: 2,
          fill: true,
          type: 'bar',
          data: this.entries
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
                hour: 'hh A'
              }
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
              max: 100,
              stepSize: 10,
              callback: function (value, index, values) {
                return value + '%'
              }
            }
          }]
        }
      }
    }
  },

  mounted () {
    this.generateLabels(this.interval)
    this.generateChart()
    this.renderChart(this.datacollection, this.options)
  }

}

</script>
