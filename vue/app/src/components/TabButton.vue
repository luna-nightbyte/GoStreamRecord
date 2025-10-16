<template>
    <button class="tab-button" @mouseover="handleMouseOver" @mouseout="handleMouseOut">
      <span class="tab-button-font">{{ displayText }}</span>
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

/* 
 TAB
*/


  </style>
  