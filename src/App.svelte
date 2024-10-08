<script>
  /* eslint-env browser */
  const OBS_WEBSOCKET_LATEST_VERSION = '5.0.1' // https://api.github.com/repos/Palakis/obs-websocket/releases/latest

  // Imports
  import { onMount } from 'svelte'
  import {
    mdiRecord,
    mdiStop,
    mdiPause,
    mdiPlayPause,
    mdiConnection,
  } from '@mdi/js'
  import Icon from 'mdi-svelte'
  import { compareVersions } from 'compare-versions'

  import './style.scss'
  import { obs, sendCommand } from './obs.js'
  import ProgramPreview from './ProgramPreview.svelte'

  onMount(async () => {
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.register('/service-worker.js')
    }

    // Request screen wakelock
    if ('wakeLock' in navigator) {
      try {
        await navigator.wakeLock.request('screen')
        // Re-request when coming back
        document.addEventListener('visibilitychange', async () => {
          if (document.visibilityState === 'visible') {
            await navigator.wakeLock.request('screen')
          }
        })
      } catch (e) {}
    }

    // Toggle the navigation hamburger menu on mobile
    const navbar = document.querySelector('.navbar-burger')
    navbar.addEventListener('click', () => {
      navbar.classList.toggle('is-active')
      document
        .getElementById(navbar.dataset.target)
        .classList.toggle('is-active')
    })

    if (document.location.hash !== '') {
      // Read address from hash
      address = document.location.hash.slice(1)

      // This allows you to add a password in the URL like this:
      // http://obs-web.niek.tv/#ws://localhost:4455#password
      if (address.includes('#')) {
        [address, password] = address.split('#')
      }
      await connect()
    }

    // Export the sendCommand() function to the window object
    window.sendCommand = sendCommand
  })

  // State
  let connected
  let heartbeat = {}
  let heartbeatInterval

  let isSceneOnTop = window.localStorage.getItem('isSceneOnTop') || false
  let isIconMode = window.localStorage.getItem('isIconMode') || false

  let address = "ws://0.0.0.0:4455"
  let password = "xxxxxxxxxxxxxxxx"

  let errorMessage = ''
  let imageFormat = 'jpg'

  $: isSceneOnTop
    ? window.localStorage.setItem('isSceneOnTop', 'true')
    : window.localStorage.removeItem('isSceneOnTop')

  $: isIconMode
    ? window.localStorage.setItem('isIconMode', 'true')
    : window.localStorage.removeItem('isIconMode')

  function formatTime (secs) {
    secs = Math.round(secs / 1000)
    const hours = Math.floor(secs / 3600)
    secs -= hours * 3600
    const mins = Math.floor(secs / 60)
    secs -= mins * 60
    return hours > 0
      ? `${hours}:${mins < 10 ? '0' : ''}${mins}:${secs < 10 ? '0' : ''}${secs}`
      : `${mins < 10 ? '0' : ''}${mins}:${secs < 10 ? '0' : ''}${secs}`
  }

  
  async function startRecording () {
    await sendCommand('StartRecord')
  }

  async function stopRecording () {
    await sendCommand('StopRecord')
  }


  async function pauseRecording () {
    await sendCommand('PauseRecord')
  }

  async function resumeRecording () {
    await sendCommand('ResumeRecord')
  }

  async function connect() {
    address = address || 'ws://0.0.0.0:4455';
    if (address.indexOf('://') === -1) {
        const secure = location.protocol === 'https:' || address.endsWith(':443');
        address = secure ? 'wss://' + address : 'ws://' + address;
    }
    console.log('Connecting to:', address, '- using password:', password);
    await disconnect();  // Ensure this is properly defined elsewhere
    
    try {
        const { obsWebSocketVersion, negotiatedRpcVersion } = await obs.connect(address, password);
        console.log(
            `Connected to obs-websocket version ${obsWebSocketVersion} (using RPC ${negotiatedRpcVersion})`
        );
    } catch (e) {
        console.log('Connection failed:', e.message);
        errorMessage = e.message;

    }
}

async function startOBS() {
  try {
    // Make a POST request to your backend endpoint to start OBS Studio
    const response = await fetch('http://0.0.0.0:3001/start-obs', { method: 'GET' });
    const data = await response.json();

    if (data.success) {
      console.log('OBS Studio started successfully.');
      // Optionally, attempt to connect to OBS after a delay
      await connect(); // Adjust delay as needed
    } else {
      console.error('Failed to start OBS Studio:', data.error);
    }
  } catch (error) {
    console.error('Error starting OBS Studio:', error.message);
  }
}




  async function disconnect () {
    await obs.disconnect()
    clearInterval(heartbeatInterval)
    connected = false
    errorMessage = 'Disconnected'
  }

  // OBS events
  obs.on('ConnectionClosed', () => {
    connected = false
    window.history.pushState(
      '',
      document.title,
      window.location.pathname + window.location.search
    ) // Remove the hash
    console.log('Connection closed')
  })

  obs.on('Identified', async () => {
    console.log('Connected')
    connected = true
    document.location.hash = address // For easy bookmarking
    const data = await sendCommand('GetVersion')
    const version = data.obsWebSocketVersion || ''
    console.log('OBS-websocket version:', version)
    if (compareVersions(version, OBS_WEBSOCKET_LATEST_VERSION) < 0) {
      alert(
        'You are running an outdated OBS-websocket (version ' +
          version +
          '), please upgrade to the latest version for full compatibility.'
      )
    }
    if (
      data.supportedImageFormats.includes('webp') &&
      document
        .createElement('canvas')
        .toDataURL('image/webp')
        .indexOf('data:image/webp') === 0
    ) {
      imageFormat = 'webp'
    }
    heartbeatInterval = setInterval(async () => {
      const stats = await sendCommand('GetStats')
      const streaming = await sendCommand('GetStreamStatus')
      const recording = await sendCommand('GetRecordStatus')
      heartbeat = { stats, streaming, recording }
      // console.log(heartbeat);
    }, 1000) // Heartbeat

  })

  obs.on('ConnectionError', async () => {
    errorMessage = 'Please enter your password:'
    document.getElementById('password').focus()
    if (!password) {
      connected = false
    } else {
      await connect()
    }
  })

</script>

<svelte:head>
  <title>OBS-web remote control</title>
</svelte:head>

<nav class="navbar is-primary" aria-label="main navigation">
  <div class="navbar-brand">
    <a class="navbar-item is-size-4 has-text-weight-bold" href="/">
      <img src="favicon.png" alt="OBS-web" class="rotate" /></a
    >

    <!-- svelte-ignore a11y-missing-attribute -->
    <button
      class="navbar-burger burger"
      aria-label="menu"
      aria-expanded="false"
      data-target="navmenu"
    >
      <span aria-hidden="true" />
      <span aria-hidden="true" />
      <span aria-hidden="true" />
    </button>
  </div>

  <div id="navmenu" class="navbar-menu">
    <div class="navbar-end">
      <div class="navbar-item">
      
      </div>
    </div>
  </div>
</nav>

<section class="controls-section mt-5">
  <div class="container">
    <div class="buttons is-centered">
      <!-- svelte-ignore a11y-missing-attribute -->
      {#if connected}
        {#if heartbeat && heartbeat.recording && heartbeat.recording.outputActive}
          {#if heartbeat.recording.outputPaused}
            <button
              class="button is-large is-danger"
              on:click={resumeRecording}
              title="Resume Recording"
            >
              <span class="icon"><Icon path={mdiPlayPause} /></span>
              <span class="ml-2">Resume</span>
            </button>
          {:else}
            <button
              class="button is-large is-success"
              on:click={pauseRecording}
              title="Pause Recording"
            >
              <span class="icon"><Icon path={mdiPause} /></span>
              <span class="ml-2">Pause</span>
            </button>
          {/if}
          <button
            class="button is-large is-danger"
            on:click={stopRecording}
            title="Stop Recording"
          >
            <span class="icon"><Icon path={mdiStop} /></span>
            <span class="ml-2">Stop</span>
            <span class="ml-2">{formatTime(heartbeat.recording.outputDuration)}</span>
          </button>
        {:else}
          <button
            class="button is-large is-danger is-light"
            on:click={startRecording}
            title="Start Recording"
          >
            <span class="icon"><Icon path={mdiRecord} /></span>
            <span class="ml-2">Start</span>
          </button>
        {/if}

        <button
          class="button is-large is-danger is-light"
          on:click={disconnect}
          title="Disconnect"
        >
          <span class="icon"><Icon path={mdiConnection} /></span>
          <span class="ml-2">Disconnect</span>
        </button>
      {:else}
        
      {/if}
    </div>
  </div>
</section>


<section class="section">
  <div class="container">
    
    {#if connected}

      {#if isSceneOnTop}
        <ProgramPreview {imageFormat} />
      {/if}
      {#if !isSceneOnTop}
        <ProgramPreview {imageFormat} />
      {/if}
      
    {:else}
    <div class="field has-text-centered">
      <button
      class="button is-success is-rounded is-medium"
      on:click={() => {   console.log('Start OBS clicked'); startOBS(); connect(); }}
      title="Start OBS Studio"
    >
      Start OBS
    </button></div>
        
      <h1 class="subtitle welcome-message">
        Capture every moment like a pro—remotely, effortlessly, and in style.
        Control your
        <strong>OBS Studio</strong>
        remotely. 
        
      </h1>

      {#if document.location.protocol === 'https:'}
        <div class="notification is-danger">
          You are checking this page on a secure HTTPS connection. That's great,
          but it means you can
          <strong>only</strong>
          connect to WSS (secure websocket) addresses, for example OBS exposed with
          <a href="https://ngrok.com/">ngrok</a>
          or
          <a href="https://pagekite.net/">pagekite</a>
          . If you want to connect to a local OBS instance,
          <strong>
            <a
              href="http://{document.location.hostname}{document.location.port
                ? ':' + document.location.port
                : ''}{document.location.pathname}"
            >
              please click here to load the non-secure version of this page
            </a>
          </strong>
          .
        </div>
      {/if}
    {/if}
  </div>
</section>

<footer class="footer">
  <div class="content has-text-centered">
    <p>
      <strong>OBS-web-mini</strong>
      by
      <a href="http://www.adithyasudhan.netlify.app/">Adi</a>
      &mdash; see
      <a href="https://github.com/xocolatl-aficionado/obs-web-mini/">GitHub</a>
      for source code.
    </p>
  </div>
</footer>
