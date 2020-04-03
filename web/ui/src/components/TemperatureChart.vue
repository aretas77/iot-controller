<script>
import { Line } from 'vue-chartjs'

export default {
  extends: Line,
  data () {
    return {
      datacollection: {
        labels: ['Data', 'February', 'March'],
        datasets: [
          {
            label: 'Data One',
            backgroundColor: '#f87979',
            pointBackgroundColor: 'white',
            borderWidth: 2,
            pointBorderColor: '#249EBF',
            data: [27, 24, 22, 24]
          },
          {
            label: 'Data One',
            backgroundColor: 'green',
            pointBackgroundColor: 'white',
            borderWidth: 2,
            pointBorderColor: 'yellow',
            data: [25, 26, 23, 25]
          }
        ]
      },
      // Chart.js options that control the appearance of the chart
      options: {
        scales: {
          yAxes: [{
            ticks: {
              beginAtZero: true
            },
            gridLines: {
              display: true
            }
          }],
          xAxes: [{
            gridLines: {
              display: false
            }
          }]
        },
        legend: {
          display: true
        },
        responsive: true,
        maintainAspectRatio: false
      }
    }
  },
  mounted () {
    this.renderChart(this.datacollection, this.options)
  }
}
</script>
