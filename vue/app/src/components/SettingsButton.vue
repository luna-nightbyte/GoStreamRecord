<template>
    <button class="ultra-cool-button" @mouseover="handleMouseOver" @mouseout="handleMouseOut">
      <span>{{ displayText }}</span>
      <div class="ripple"></div>
    </button>
  </template>
  
  <script>
  export default {
    props: {
      text: {
        type: String,
        required: true
      },
      hovertext: {
        type: String,
        required: true
      }
    },
    data() {
      return {
        originalText: this.text,
        displayText: this.text,
        originalHoverText: this.hovertext
      };
    },
    methods: {
      handleMouseOver() {
        this.displayText = this.originalText;
      },
      handleMouseOut() {
        this.displayText = this.originalText;
      }
    }
  };
  </script>
  
  <style scoped>
  </style>
  