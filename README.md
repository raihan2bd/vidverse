<a name="readme-top"></a>
<h1 align='center'>Vidverse</h1>


# 📗 Table of Contents

- [📗 Table of Contents](#-table-of-contents)
- [ Vidverse ](#-about-project-)
  - [🛠 Built With ](#-built-with-)
    - [Tech Stack ](#tech-stack-)
    - [Key Features ](#key-features-)
  - [💻 Getting Started ](#-getting-started-)
    - [Prerequisites](#prerequisites)
    - [Setup](#setup)
    - [Install](#install)
    - [Database](#database)
    - [Usage](#usage)
    - [Build](#build)
    - [Deployment](#deployment)
  - [👥 Authors ](#-authors-)
  - [🔭 Future Features ](#-future-features-)
  - [🤝 Contributing ](#-contributing-)
  - [⭐️ Show your support ](#️-show-your-support-)
  - [🙏 Acknowledgments ](#-acknowledgments-)
  - [📝 License ](#-license-)


# Vidverse <a name="about-project"></a>
Introducing Vidverse – a dynamic, full-stack web application crafted with Golang (Go), Next.js, TypeScript, PostgreSQL, and a host of other cutting-edge technologies. Vidverse offers users a seamless experience to watch, like, comment, and engage with an eclectic array of content, and users can also authenticate using (credentials and social platforms).

Creators have the power to effortlessly share their latest videos and manage their content, while subscribers revel in curated feeds, saved favorites, and personalized watch histories. User can also see their notifications.

This is the back-end version. If you want to see the front-end part please visit [here](https://github.com/raihan2bd/vidverse-client)

## 🛠 Built With <a name="built-with"></a>
### Tech Stack <a name="tech-stack"></a>

<details open>
  <summary>Front End</summary>
  <ul>
    <li>Nextjs</li>
    <li>React</li>
    <li>TypeScript</li>
    <li>Html</li>
    <li>CSS</li>
    <li>Tailwind CSS</li>
  </ul>
</details>
<details open>
  <summary>Back End</summary>
  <ul>
    <li>Golang</li>
    <li>PostgreSQL</li>
    <li>Gin</li>
    <li>Gorm</li>
    <li>Web Socket</li>
    <li>Firebase</li>
  </ul>
</details>


<!-- LIVE DEMO -->

## 🚀 Live Demo <a name="live-demo"></a>

- Project ScreenShot

***Update soon***


<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Key Features <a name="key-features"></a>

- Users can authenticate using their credentials or through various social platforms, ensuring a flexible and convenient login experience.
- Creators can effortlessly share their latest videos and manage their content with ease, fostering a vibrant community of content producers.
- Subscribers enjoy curated feeds, saved favorites, and personalized watch histories, enhancing their viewing experience.
- Users receive notifications, ensuring they stay updated on the latest activities and interactions within the Vidverse community.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 💻 Getting Started <a name="getting-started"></a> 

To get a local copy up and running, follow these steps.

### Prerequisites

In order to run this project you need:
- Then Make sure you have installed [Go (golang)](https://go.dev/dl/) version 1.20.4 or the latest stable version.
- Then make sure you have installed [PostgreSQL](https://www.postgresql.org/) on your local machine if you want to use this project locally.
- Then Create a database called `vidverse`

- First of all to see this project's graphical interface make sure you run the [front-end](https://github.com/raihan2bd/filmwise-front) part

### Setup

- Clone this repository to your desired folder:

```sh
  cd your-folder
  https://github.com/raihan2bd/filmwise.git
```

- Before running the project please make sure you create a `.env` or rename the `.env.example` file to `.env` file and update the following credentials
  ```
  PORT=Your port number
  DB_URI="host=localhost user=postgres password=#Dovaridu90 dbname=vidverse port=5432 sslmode=disable"
  SECRET=yoursecret
  CLD_URI= your cloudinary uri
  CLOUD_NAME=your cloud name
  JWT_SECRET=your jwt secret
  ENVIRONMENT= your environment
  
  MAIL_DOMAIN= your mail domain
  MAIL_HOST= your mail host
  MAIL_PORT= your mail port
  MAIL_USERNAME= your mail username
  MAIL_PASSWORD= your mail password
  MAIL_ENCRYPTION= your mail encryption
  MAIL_FROM_NAME= your mail from name
  MAIL_FROM_ADDRESS= your mail from address
  ```

- Then don't forget to create a service account to implement the social login. Then add your service account credentials inside a file called `service-account.json` make sure you add this file in your project root. see the [Blog](https://sharma-vikashkr.medium.com/firebase-how-to-setup-a-firebase-service-account-836a70bb6646) if you need more information about it.

### Getting JWT Secret Key

To obtain the JWT secret key, please <a href="https://go.dev/play/" target="_blank" style="color: blue; font-size: 12px; text-decoration: underline;">click here</a>.

An open Go terminal is required. In the terminal, please copy the following code and paste it into your <span style="color: lightblue;">.env</span> file.

Like this

JWT_SECRET="*******"

```

package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	// Define the desired length of the secret key in bytes
	keyLength := 64 // Adjust the length as needed

	// Create a byte slice to hold the random bytes
	key := make([]byte, keyLength)

	// Generate random bytes using crypto/rand
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("Error generating random bytes:", err)
		return
	}

	// Encode the random bytes in base64 to create a string
	secretKey := base64.StdEncoding.EncodeToString(key)

	fmt.Println("Generated JWT secret key:", secretKey)
}

```

### Getting Cloudinary Serect key and Name

Cloudinary is a cloud-based media management platform that helps businesses and developers efficiently store, manage, and deliver images and videos for websites and applications. It provides features like image and video uploading, storage, transformation, optimization, and content delivery via a content delivery network (CDN), making it easier to handle media assets in web and mobile applications. Cloudinary's services can enhance website performance, user experience, and streamline media asset workflows.

For using this you need to <a href="https://cloudinary.com/users/register_free" target="_blank" style="hover:text-decoration: underline; font-size:1.1rem">Create an account</a> Or if you have an Account you need to <a href="https://cloudinary.com/users/login" target="_blank" style="hover:text-decoration: underline; font-size:1.2rem;">Sign In </a>

After that go into the <span style="color: #3448c5;">Dashborad</span>

Copy the <span style="color: #3448c5;">Cloud Name</span> and
<span style="color: #3448c5;">API Environment variable</span>

```
CLD_URI="cloudinary://******"

CLOUD_NAME="****"

```

### Install

Install this project with:

- Install the required gems with:

```sh
go mod tidy
```

### Database

- Create the databases properly, You need to open an SQL editor and run the `/database/schema.sql` file script. Make sure you run the script block by block.

### Usage

- To run the development server, execute the following command:

```sh
go run ./cmd/api/ .
```

### Build

- To build the project for production-ready run the following command:

```sh
go build -o main ./cmd/api/*.go
```


### Deployment

To deploy your project online You can visit [Render](https://www.render.com/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## 👥 Author <a name="author"></a>

👤 **Abu Raihan**

- GitHub: [@raihan2bd](https://github.com/raihan2bd)
- Twitter: [@raihan2bd](https://twitter.com/raihan2bd)
- LinkedIn: [raihan2bd](https://linkedin.com/in/raihan2bd)


👤 **Nurgul Kereikhan**

- GitHub: [@githubhandle](https://github.com/NurkaAmre)
- Twitter: [@twitterhandle](https://twitter.com/AmreNurgul)
- LinkedIn: [LinkedIn](https://www.linkedin.com/in/amre-nurgul/)

  
<p align="right">(<a href="#readme-top">back to top</a>)</p>


## 🔭 Future Features <a name="future-features"></a>

- [ ] **Implement improvements to provide users with an even smoother and more enjoyable experience.**
- [ ] **Transition the application to a microservices architecture for improved scalability and maintainability.**
- [ ] **Add a new feature called `Shorts` to the platform, enabling users to create and share short-form video content.**
- [ ] **Incorporate FFmpeg to enhance video streaming capabilities similar to YouTube, and leverage AWS for hosting these files to ensure seamless playback and scalability.**

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## 🤝 Contributing <a name="contributing"></a>

Contributions, issues, and feature requests are welcome!

Feel free to check the [issues page](https://github.com/raihan2bd/vidverse/issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## ⭐️ Show your support <a name="support"></a>

If you like this project, please leave a ⭐️

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## 🙏 Acknowledgments <a name="acknowledgements"></a>

We extend our heartfelt gratitude to [Microverse](https://microverse.org) and [Trevor Sawler](https://www.gocode.ca/) for their invaluable assistance in mastering the tech stack utilized in this project. Additionally, we express our sincere appreciation to [Cloudinary](https://cloudinary.com/) for generously providing us with free cloud space.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## 📝 License <a name="license"></a>

This project is [MIT](./LICENSE) licensed.

<p align="right">(<a href="#readme-top">back to top</a>)</p>
