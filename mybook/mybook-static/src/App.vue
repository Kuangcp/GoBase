<template>
  <div id="app">
    <MainTabs/>
  </div>
</template>

<script>
import MainTabs from "./page/MainTabs";

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
  /*https://pxlab.cn/color/index.html*/
  /*background-image: linear-gradient(25deg, #65b8b5, #96b4ae, #bbb0a7, #dcaaa0);*/
  /*background-image: linear-gradient(25deg, #227495, #76859d, #b095a5, #e7a5ad)*/
  background: linear-gradient(-180deg, #BCC5CE 0%, #929EAD 98%), radial-gradient(at top left, rgba(255,255,255,0.30) 0%, rgba(0,0,0,0.30) 100%); background-blend-mode: screen;
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
