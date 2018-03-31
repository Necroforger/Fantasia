<script>
import { Line } from "vue-chartjs";

export default {
  name: "LineChart",
  extends: Line,
  props: ["label", "autostart", "endpoint", "max", "automax"],
  data() {
    return {
      context: {
        labels: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
        datasets: [
          {
            pointRadius: 0,
            tension: 0,
            label: this.label,
            borderWidth: 0.7,
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
                    beginAtZero: !this.autostart
                  }
                : {
                    beginAtZero: !this.autostart,
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
    }, 5000);
  },
  methods: {
    // update the chart
    fetchUpdate() {
      console.log(`${this.label}: attempting to fetch updates`);
      if (!this.endpoint) {
        console.log(`${this.label}: endpoint not set. update failed`);
        return;
      }
      let xhr = new XMLHttpRequest();
      xhr.onreadystatechange = () => {
        if (xhr.readyState == 4) {
          if (xhr.status != 200) {
            console.log(
              `${this.label}: request did not return 200 OK [${xhr.status}]`
            );
            return;
          }
          if (!xhr.response) {
            console.log(`${this.label}: xhr.response does not exist`);
            return;
          }
          console.log(`${this.label}: updating with:` + xhr.response);

          this.updatePoints(xhr.response.data);
          this.render();
        }
      };
      xhr.open("GET", this.endpoint, true);
      xhr.responseType = "json";
      xhr.send();
    },

    updatePoints(dataset) {
      for (i in dataset.data) {
        this.context.datasets[0].data = dataset.data[i];
        this.context.labels = dataset.labels[i];
      }
    },

    // render the chart
    render() {
      this.renderChart(this.context, this.config);
    }
  }
};
</script>
