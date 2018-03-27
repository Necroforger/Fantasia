<script>
import { Line } from "vue-chartjs";

export default {
  name: "LineChart",
  extends: Line,
  props: ["label", "endpoint", "max", "automax"],
  data() {
    return {
      context: {
        labels: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
        datasets: [
          {
            label: this.label,
            backgroundColor: "rgba(255, 0, 255, 0.09)",
            borderColor: "green",
            data: (() => {
              let retval = [];
              for (let i = 0; i < 10; i++) {
                retval.push(Math.random() * 100);
              }
              return retval;
            })()
          }
        ]
      },
      config: {
        maintainAspectRatio: false,
        responsive: true,
        animation: false,
        scales: {
          yAxes: [
            {
              ticks: this.automax
                ? {
                    beginAtZero: true
                  }
                : {
                    beginAtZero: true,
                    max: this.max ? parseInt(this.max) : 100
                  }
            }
          ]
        }
      }
    };
  },
  mounted() {
    // Overwriting base render method with actual data.
    this.renderChart(this.context, this.config);

    this.fetchUpdate();
    setInterval(() => {
      this.fetchUpdate();
    }, 1000);
  },
  methods: {
    fetchUpdate() {
      console.log(`${this.label}: attempting to fetch updates`);
      if (!this.endpoint) {
        console.log(`${this.label}: endpoint not set. update failed`);
        return;
      }
      let xhr = new XMLHttpRequest();
      xhr.onreadystatechange = () => {
        if (xhr.readyState == 4) {
          console.log(`${this.label}: updating with:` + xhr.response);
          this.context.datasets[0].data = xhr.response.data;
          this.context.labels = xhr.response.labels;
          this.render();
        }
      };
      xhr.open("GET", this.endpoint, true);
      xhr.send();
    },
    render() {
      this.renderChart(this.context, this.config);
    }
  }
};
</script>
