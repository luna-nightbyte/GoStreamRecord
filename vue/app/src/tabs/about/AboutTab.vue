<template>
  <div class="item-container ">
    <h1>About Our Video Downloader Tool</h1>
    <p>
      Welcome to our Video Downloader Tool! We are thrilled to offer a simple, efficient,
      and completely free solution for downloading videos from the web. Whether you are
      looking to save a tutorial, a music video, or any other type of online content, our
      tool provides a seamless experience to help you download videos directly to your device.
    </p>

    <h2>Key Features</h2>
    <ul>
      <li><strong>Record live streams:</strong> Record and manage livestreams directly in the webUI.</li>
       <li><strong>Download Videos Directly from URL:</strong> Easily download any video by simply pasting the URL into our application. No complex steps, just straightforward downloads.</li>
      <li><strong>Highest Quality:</strong> The tool always looks for the highest quality available for the selected videos.  </li>
      <li><strong>Batch Download:</strong> Save time by downloading an entire page of videos in one go. Simply search for your desired videos, select the page, and initiate the batch download.</li>
      <li><strong>Privacy First:</strong> We respect your privacy. Our application does not connect to any external servers or collect any user data. Your download activities remain completely private and secure.</li>
      <li><strong>Free Forever:</strong> Enjoy all our features at no cost. Our application is free of charge, now and always.</li>
      <li><strong>Local Webserver:</strong> Our application runs on a local webserver on port 8080, ensuring fast and reliable performance without the need for an internet connection.</li>
      <li><strong>Media View:</strong> Easily manage and view all your downloaded videos within our intuitive media view interface.</li>
    </ul>

    <h2>Privacy and Security</h2>
    <p>
      We understand the importance of privacy and security in today's digital age. Our Video
      Downloader Tool is designed with these principles at its core. 
      <span class="highlight">We do not track, save, or use any user information in any way.</span>
      Your privacy is our top priority, and you can trust that your data remains completely
      anonymous and secure.
    </p>

    <h2>No External Connections</h2>
    <p>
      Another key feature of our tool is that it operates independently of any external servers.
      This means that 
      <span class="highlight">our tool is not connected to any external servers in any way.</span>
      All video processing is done locally on your device, ensuring that your downloads are fast
      and reliable without any risks of data breaches or external monitoring.
    </p>

    <h2>Completely Free to Use</h2>
    <p>
      Our mission is to provide a high-quality video downloading service to everyone, without any
      cost barriers. We are proud to offer this tool 
      <span class="highlight">completely free of charge.</span>
      There are no hidden fees, no subscriptions, and no advertisements. Simply use the tool as
      needed and enjoy the convenience it brings.
    </p>

    <h2>How It Works</h2>
    <p>Using our Video Downloader Tool is straightforward and user-friendly. Just follow these simple steps:</p>
    <ol>
      <li>Copy the URL of the video you wish to download.</li>
      <li>Paste the URL into the designated field on our tool.</li>
      <li>Click the download button and choose your preferred video format and quality.</li>
      <li>Wait for the download to complete, and enjoy your video offline!</li>
      <li>Bulk option is for downloading a page of videos. i.e "https://somesite.com/search=hello&page=3"</li>
    </ol>

    <h2>Known Issues</h2>
    <ul>
      <li><strong>Video Naming:</strong> Currently, videos do not save with their original names. Users need to provide a fixed name for each download. We are working on an update to address this issue in the future.</li>
      <li><strong>Batch Downloads:</strong> While our batch download feature is powerful, it does not work on all websites. We are continuously working to improve compatibility with more sites.</li>
    </ul>

    <p>Thank you for choosing our tool. Happy downloading!</p>
  </div>
</template>

<script>
export default {
  name: "AboutUs"
};
</script>

<style scoped>
body {
  font-family: Arial, sans-serif;
  background-color: #f4f4f4;
  color: #333;
  margin: 0;
  padding: 0;
}

.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background-color: #000000;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

h1 {
  color: #0053a7;
  text-align: center;
}

h2, h3 {
  color: #003061;
}

p {
  line-height: 1.6;
  margin: 20px 0;
}

ul {
  list-style-type: disc;
  margin-left: 20px;
}

ul li {
  margin-bottom: 10px;
}

ol {
  list-style-type: decimal;
  margin-left: 20px;
}

ol li {
  margin-bottom: 10px;
}

.highlight {
  color: #00baf8;
  font-weight: bold;
}
</style>
