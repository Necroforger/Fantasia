<template>
  <div id="app">
    <template v-if="tab==='dashboard'">
      <Dashboard></Dashboard>
    </template>
  </div>
</template>

<script>
import Dashboard from "./tabs/Dashboard";

export default {
  name: "app",
  components: {
    Dashboard,
  },
  data() {
    return {
      tab: "dashboard",
    };
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
