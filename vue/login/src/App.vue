<!-- File: src/components/Login.vue -->
<template>
  <main class="login-page">
    <section class="card" role="form" aria-labelledby="login-title">
      <header class="card__header">
        <div class="card__header-row">
          <div class="card__titles">
            <h1 id="login-title" class="title">Sign in</h1> 
          </div>
          <div class="status" :class="statusClass" aria-live="polite" :aria-label="`Status: ${statusText}`">
            <span class="status__dot" aria-hidden="true"></span>
            <span class="status__text">{{ statusText }}</span>
          </div>
        </div>
      </header>

      <form class="form" @submit.prevent="submitForm" novalidate>
        <div class="field">
          <label class="label" for="username">Username</label>
          <input
            id="username"
            class="input"
            type="text"
            v-model.trim="username"
            :disabled="loading"
            autocomplete="username"
            required
          />
        </div>

        <div class="field">
          <label class="label" for="password">Password</label>
          <div class="input input--with-btn">
            <input
              id="password"
              :type="showPassword ? 'text' : 'password'"
              v-model="password"
              :disabled="loading"
              autocomplete="current-password"
              required
              class="input__control"
            />
            <button
              class="icon-btn"
              type="button"
              :aria-pressed="showPassword ? 'true' : 'false'"
              :title="showPassword ? 'Hide password' : 'Show password'"
              @click="showPassword = !showPassword"
            >
              <span v-if="showPassword" aria-hidden="true">🙈</span>
              <span v-else aria-hidden="true">👁️</span>
              <span class="sr-only">{{ showPassword ? 'Hide password' : 'Show password' }}</span>
            </button>
          </div>
        </div>

        <button class="btn" :disabled="loading">
          <span v-if="!loading">Continue</span>
          <span v-else class="spinner" aria-hidden="true"></span>
          <span class="sr-only" v-if="loading">Signing in…</span>
        </button>
      </form>

      <!-- <footer class="card__footer">
        <a href="#" class="link" @click.prevent>Forgot password?</a>
      </footer> -->
    </section>

    <transition name="toast">
      <output
        v-if="responseMessage.message"
        class="toast"
        :class="responseMessage.isError ? 'toast--danger' : 'toast--info'"
        role="status"
        aria-live="polite"
      >
        {{ responseMessage.message }}
      </output>
    </transition>
  </main>
</template>

<script setup>
import { ref, computed } from 'vue'

const username = ref('')
const password = ref('')
const statusText = ref('Online')
const responseMessage = ref({ message: '', isError: false })
const loading = ref(false)
const showPassword = ref(false)

const statusClass = computed(() => ({
  'status--ok': statusText.value?.toLowerCase() === 'online',
  'status--warn': statusText.value?.toLowerCase() === 'degraded',
  'status--down': statusText.value?.toLowerCase() === 'offline'
}))

function showResponse(message, isError = false, timeoutMs = 4000) {
  responseMessage.value = { message, isError }
  if (timeoutMs > 0) {
    window.setTimeout(() => {
      responseMessage.value = { message: '', isError: false }
    }, timeoutMs)
  }
}

async function submitForm() {
  if (!username.value || !password.value) {
    showResponse('Please enter both username and password.', true)
    return
  }
  loading.value = true
  try {
    const formData = new FormData()
    formData.append('username', username.value)
    formData.append('password', password.value)

    const response = await fetch('/login', { method: 'POST', body: formData })

    if (response.redirected) {
      window.location.href = response.url
      return
    }

    const data = await response.json().catch(() => ({}))
    if (data && data.message) {
      showResponse(data.message, true)
    } else {
      showResponse('Unexpected response from server', true)
    }
  } catch (err) {
    console.error('Error during login:', err)
    showResponse('An unexpected error occurred', true)
  } finally {
    loading.value = false
  }
}
</script>
 