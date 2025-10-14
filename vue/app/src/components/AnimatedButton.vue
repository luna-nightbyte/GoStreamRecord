<template>
  <button class="animated-button" @mouseover="handleMouseOver" @mouseout="handleMouseOut">
    <span class="button-font">{{ displayText }}</span>
    <div class="ripple"></div>
  </button>
</template>

<script>
export default {
  props: {
    text: { type: String, required: true },
    hovertext: { type: String, required: true }
  },
  data() {
    return {
      originalText: this.text,
      displayText: this.text
    };
  },
  methods: {
    handleMouseOver() { this.displayText = this.originalText; },
    handleMouseOut() { this.displayText = this.originalText; }
  }
};
</script>

<style>
.animated-button {
  position: relative;
  display: inline-block;
  padding: 15px 30px;
  font-size: 18px;
  font-weight: bold;
  color: #fff;
  background: linear-gradient(45deg, var(--secondary), var(--primary-light));
  border: none;
  border-radius: var(--radius);
  cursor: pointer;
  overflow: hidden;
  transition: transform 0.4s, box-shadow 0.4s, background 0.4s;
  z-index: 1;
  box-shadow: var(--shadow);
}

.animated-button:hover {
  transform: scale(1.05);
  background: linear-gradient(45deg, var(--primary-light), var(--primary));
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
}

.animated-button span {
  position: relative;
  z-index: 2;
  font-family: var(--font-family);
}

.animated-button::before,
.animated-button::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 300%;
  height: 300%;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  transition: transform 0.6s;
  transform: translate(-50%, -50%) scale(0);
  z-index: 1;
}

.animated-button:hover::before,
.animated-button:hover::after {
  transform: translate(-50%, -50%) scale(1);
}

.animated-button .ripple {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 300%;
  height: 300%;
  background: rgba(255, 255, 255, 0);
  border-radius: 50%;
  transform: translate(-50%, -50%) scale(0);
  transition: transform 0.8s;
}

.animated-button:hover .ripple {
  transform: translate(-50%, -50%) scale(1);
}

.button-font {
  margin: 0;
  color: #fff;
}
</style>
