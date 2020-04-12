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
      minTemp: -20,
      datacollection: {
        // Could be: 6hrs by hour, 1hr by 60min
        labels: [],
        datasets: [
          {
            label: 'Temperature received',
            backgroundColor: '#f87979',
            pointBackgroundColor: 'white',
            borderWidth: 2,
            fill: false,
            pointBorderColor: '#5cb85c',
            type: 'line',
            lineTension: 0,
            data: []
          },
          {
            label: 'Temperature (Hades)',
            backgroundColor: 'green',
            pointBackgroundColor: 'white',
            borderWidth: 2,
            pointBorderColor: 'yellow',
            data: []
          }
        ]
      },
      // Chart.js options that control the appearance of the chart
      options: {
        scales: {
          yAxes: [{ }],
          xAxes: [{ }]
        },
        legend: {
          display: true
        },
        responsive: true,
        maintainAspectRatio: false
      }
    }
  },

  methods: {
    // 2020-04-03T20:15:20Z
    parseData (timestamp) {
      // can select which hour
      if (timestamp === 'hr') {
        // show data in 1hr - 1min step
        this.labelsData = this.entries.map(data => {
          return this.$moment(data.temp_read_time, this.$moment.ISO_8601).format('m')
        })
      } else if (timestamp === '6hr') {
        // show data in 6hr - 6min step
      } else if (timestamp === '24hr') {
        // show current day data
        this.constructedData = this.entries.map(entry => {
          // Check if current day is today
          if (this.$moment(entry.temp_read_time).isSame(this.$moment(), 'day')) {
            return {
              x: this.$moment(entry.temp_read_time, this.$moment.ISO_8601).format('HH:mm'),
              y: entry.temperature
            }
          }
        })
        // console.log(this.constructedData)
      }
      this.datacollection.datasets[0].data = this.constructedData

      // Construct scales
      this.options.scales.yAxes = [{
        ticks: {
          min: this.minTemp,
          max: this.maxTemp
        },
        gridLines: {
          display: true
        }
      }]

      this.options.scales.xAxes = [{
        ticks: {
          source: 'data',
          major: {
            enabled: true
          },
          maxTicksLimit: 10,
          autoSkip: true
        },
        gridLines: {
          display: true
        },
        type: 'time',
        distribution: 'series',
        offset: true,
        time: {
          unit: 'minute',
          parser: 'HH:mm',
          displayFormats: {
            minute: 'HH:mm'
          }
        }
      }]
    },
    // should be called after data is generated
    generateLabels24hr () {
    }
  },

  mounted () {
    this.parseData('24hr')

    // console.log(this.datacollection)
    this.renderChart(this.datacollection, this.options)
  }
}
</script>
