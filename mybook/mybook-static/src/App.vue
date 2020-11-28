<template>
  <div id="app">
    <Tabs/>
  </div>
</template>

<script>
import Tabs from './page/Tabs.vue'


export default {
  name: 'App',
  components: {
    Tabs
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /*text-align: center;*/
  color: #2c3e50;
  margin-top: 60px;
}
</style>
