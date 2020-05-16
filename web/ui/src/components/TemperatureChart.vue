<script>
import { Line } from 'vue-chartjs'

export default {
  extends: Line,

  props: {
    entries: {
      type: Array,
      required: true
    },
    labels: {
      type: Array,
      required: false
    }
  },

  data () {
    return {
      constructedData: [],
      maxTemp: 40,
      minTemp: -10,
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
          label: 'Temperature',
          pointBackgroundColor: 'white',
          borderWidth: 2,
          fill: false,
          // pointBorderColor: '#5cb85c',
          pointBorderColor: function (context) {
            var index = context.dataIndex
            var dataset = context.dataset
            if (index === dataset.data.length - 1) {
              return 'red'
            } else {
              return '#5cb85c'
            }
          },
          type: 'line',
          lineTension: 0,
          data: this.entries
        }
      ]

      this.options = {
        legend: {
          display: true
        },
        scales: {
          xAxes: [{
            ticks: {
              source: 'data',
              major: {
                enabled: true
              },
              maxTicksLimit: 20,
              autoSkip: true
            },
            gridLines: {
              display: true
            },
            type: 'time',
            distribution: 'series',
            offset: true,
            time: {
              unit: 'hour',
              parser: 'HH:mm',
              displayFormats: {
                minute: 'HH:mm',
                hour: 'HH mm'
              }
            }
          }],
          yAxes: [{
            ticks: {
              stepSize: 5,
              min: this.minTemp,
              max: this.maxTemp,
              callback: function (value, index, values) {
                return value + '\xB0C'
              }
            },
            gridLines: {
              display: true
            }
          }]
        },
        responsive: true,
        maintainAspectRatio: false
      }
    }
    // 2020-04-03T20:15:20Z
    // parseData (timestamp) {
    //   // can select which hour
    //   if (timestamp === 'hr') {
    //     // show data in 1hr - 1min step
    //     this.labelsData = this.entries.map(data => {
    //       return this.$moment(data.temp_read_time, this.$moment.ISO_8601).format('m')
    //     })
    //   } else if (timestamp === '24hr') {
    //     // show current day data
    //     this.constructedData = this.entries.map(entry => {
    //       // Check if current day is today
    //       if (this.$moment(entry.temp_read_time).isSame(this.$moment(), 'day')) {
    //         return {
    //           x: this.$moment(entry.temp_read_time, this.$moment.ISO_8601).format('HH:mm'),
    //           y: entry.temperature
    //         }
    //       }
    //     })
    //     console.log(this.constructedData)
    //   }
    //   this.datacollection.datasets[0].data = this.constructedData

    //   // Construct scales
    //   this.options.scales.yAxes = [{
    //     ticks: {
    //       stepSize: 5,
    //       min: this.minTemp,
    //       max: this.maxTemp,
    //       callback: function (value, index, values) {
    //         return value + '\xB0C'
    //       }
    //     },
    //     gridLines: {
    //       display: true
    //     }
    //   }]

    //   this.options.scales.xAxes = [{
    //     ticks: {
    //       source: 'data',
    //       major: {
    //         enabled: true
    //       },
    //       maxTicksLimit: 10,
    //       autoSkip: true
    //     },
    //     gridLines: {
    //       display: true
    //     },
    //     type: 'time',
    //     distribution: 'series',
    //     offset: true,
    //     time: {
    //       unit: 'minute',
    //       parser: 'HH:mm',
    //       displayFormats: {
    //         minute: 'HH:mm'
    //       }
    //     }
    //   }]
    // }
  },
  watch: {
    entries () {
      this.renderChart(this.datacollection, this.options)
    }
  },
  mounted () {
    // this.parseData('24hr')
    this.generateLabels(60)
    this.generateChart()
    this.renderChart(this.datacollection, this.options)
  }
}
</script>
