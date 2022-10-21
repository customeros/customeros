<div align="center">
  <a href="https://openline.ai">
    <img
      src="https://www.openline.ai/TeamHero.svg"
      alt="Openline Logo"
      height="64"
    />
  </a>
  <br />
  <p>
    <h3>
      <b>
        customerOS analytics API
      </b>
    </h3>
  </p>
  <p>
    GraphQL APIs to browse customer analytics data
  </p>
  <p>

[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen?logo=github)][customerOS-repo] 
[![license](https://img.shields.io/badge/license-Apache%202-blue)][apache2] 
[![stars](https://img.shields.io/github/stars/openline-ai/openline-customer-os?style=social)][customerOS-repo] 
[![twitter](https://img.shields.io/twitter/follow/openlineAI?style=social)][twitter] 
[![slack](https://img.shields.io/badge/slack-community-blueviolet.svg?logo=slack)][slack]

  </p>
  <p>
    <sub>
      Built with ❤︎ by the
      <a href="https://openline.ai">
        Openline
      </a>
      community!
    </sub>
  </p>
</div>


## 🤝 Dependencies

1.  Install [gorm][gorm]

    ```
    go get -u gorm.io/gorm
    ```
      
2. Install [gqlgen][gqlgen]
   
   *Initialize a new go module*

   ```
   mkdir example
   cd example
   go mod init example
   ```
   
   *Add `github.com/99designs/gqlgen` to your project’s `tools.go`*
   
   ```
   printf '// +build tools\npackage tools\nimport (_ "github.com/99designs/gqlgen"\n _ "github.com/99designs/gqlgen/graphql/introspection")' | gofmt > tools.go
   ```

## 🚀 Quick start


1. Add any missing module requirements necessary to build the current module’s packages and dependencies, and remove requirements on modules that don’t provide any relevant packages.

       go mod tidy

2. Generate graphql models

       go generate ./...

3. Set environment variables for DB connection:
   1. DB_HOST
   2. DB_PORT
   3. DB_NAME
   4. DB_USER
   5. DB_PWD


4. Start the graphql server

       go run server.go
       
## 💪 Contributions

We ❤️ contributions!  If you'd like to help out, check out our [contributors guide][contributions].
       
[apache2]: https://www.apache.org/licenses/LICENSE-2.0
[contributions]: https://github.com/openline-ai/community/blob/main/README.md
[customerOS-repo]: https://github.com/openline-ai/openline-customer-os/
[gorm]: https://github.com/go-gorm/gorm
[gqlgen]: https://github.com/99designs/gqlgen
[slack]: https://join.slack.com/t/openline-ai/shared_invite/zt-1i6umaw6c-aaap4VwvGHeoJ1zz~ngCKQ
[twitter]: https://twitter.com/OpenlineAI
