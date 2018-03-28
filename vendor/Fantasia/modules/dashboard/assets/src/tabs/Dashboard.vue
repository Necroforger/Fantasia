<template>
  <div class="row-three">
    <!-- CPU Usage -->
    <Card title="CPU usage" content=''>
      <LineChart endpoint="/api/stats/cpu/" label="CPU Usage" class="half-height" id="cpu-usage"/>
    </Card>

    <!-- Memory usage -->
    <Card title="Memory usage">
      <LineChart endpoint="/api/stats/mem/" label="Memory usage" class="half-height" id="memory-usage"/>
    </Card>

    <!-- Network download -->
    <Card title="Download">
      <LineChart automax="true" endpoint="/api/stats/download/" label="BytesRecv" class="half-height"></LineChart>
    </Card>

    <!-- Network upload -->
    <Card title="Upload">
      <LineChart automax="true" endpoint="/api/stats/upload/" label="BytesSent" class="half-height"></LineChart>
    </Card>
  </div>
</template>

<script>
import LineChart from "../components/LineChart";
import Card from "../components/Card";

export default {
  components: {
    LineChart,
    Card
  },
  name: "Dashboard",
  data() {
    return {};
  }
};
</script>

<style>
body {
  margin: 0px;
  background-color: #232323;
}

.row-two {
  display: grid;
  grid-template-columns: repeat(2, calc(100% / 2));
  /* grid-column-gap: 50px; */
}

.row-three {
  display: grid;
  grid-template-columns: repeat(3, calc(100% / 3));
}

.half-height {
  height: 300px;
}

#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #939393;
}

a {
  color: #42b983;
}
</style>