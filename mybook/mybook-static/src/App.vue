<template>
  <div id="app">
    <MainTabs/>
  </div>
</template>

<script>
import MainTabs from "./page/MainTabs.vue";

export default {
  name: "App",
  components: {
    MainTabs,
  },
};
</script>

<style>
* {
  margin: 0;
  padding: 0;
}

body {
  height: 100vh;
  background-image: linear-gradient(25deg, #65b8b5, #96b4ae, #bbb0a7, #dcaaa0);
  /*background-image: linear-gradient(25deg, #227495, #76859d, #b095a5, #e7a5ad)*/
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /* text-align: center; */
  color: #2c3e50;
  padding-top: 8px;
}
</style>
