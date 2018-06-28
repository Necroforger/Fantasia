<template>
  <div class="stat">
    <label class="title">{{title}}</label>
    <p class="content">{{content}}</p>
  </div>
</template>

<script>
import Card from "./Card";

export default {
  name: "Stat",
  components: { Card },
  props: ["endpoint", "title"],
  data() {
    return {
      content: "",
    };
  },
  mounted() {
    this.update();
    setInterval(this.update, 6000);
  },
  methods: {
    update() {
      var xhr = new XMLHttpRequest();
      xhr.onreadystatechange = () => {
        if (xhr.readyState == 4) {
          if (xhr.status != 200 || !xhr.response) {
            this.log("error retrieving data");
            return;
          }
          this.content = xhr.response.content;
        }
      };
      xhr.open("GET", this.endpoint, true);
      xhr.responseType = "json";
      xhr.send();
    },
    log(txt) {
      console.log(this.title + ": " + txt);
    }
  }
};
</script>

<style>
.stat {
  overflow: hidden;
  padding: 10px;
  margin: 2px;
  border: 1px dashed #303030;
  background-color: #1e1e1e;
  border-width: 2px;
  border-color: #303030;
}
.stat .title {
  color: rgb(197, 27, 197);
  font-weight: 300;
  font-size: 150%;
  text-transform: capitalize;
}
.stat .content {
  color: green;
  display: block;
  font-size: 150%;
  padding: 0px;
  margin: 0px 0px 0px 10px;
}
</style>