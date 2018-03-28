<template>
  <div class="card">
    <h1>{{title}}</h1>
    <span>{{content}}</span>
    <div>
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: "Card",
  props: ["title", "content"],
  data() {
    return {};
  }
};
</script>

<style>
.card {
  border: 5px solid #303030;
  padding-left: 10px;
  padding-right: 10px;
  border-style: dashed;
  margin: 10px;
  background-color: #202020;
}
.card h1 {
  font-size: 20px;
  margin-bottom: 1px;
  margin-top: 5px;
}
.card hr {
  color: #303030;
}
</style>