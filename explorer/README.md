RETS Explorer 
======

RETS client with a ReactJS front end, Gorilla RPC services on the backend.  

Find me at gophers.slack.com#gorets

Simplest way to run the explorer
```
docker run --rm -it -p 8080:8080 docker.io/jpfielding/goretsexplorer:latest
```

[Example Explorer UI](cmds/explorer/main.go)
```
// All hosted from go
cd explorer/client
npm i
npm run build
go run ../../cmds/explorer/main.go -port 8000 -react ./build

// To run in dev mode
npm i
cd explorer/client
npm run start
go run ../../cmds/explorer/main.go -port 8080


```
## Technologies

The front end build process uses:

### Building

- [Eslint](http://eslint.org/)
	- [airbnb preset](https://github.com/airbnb/javascript)
- [Webpack](https://webpack.github.io/)

### Javascript

- [Babel](http://babeljs.io/)
- [React](https://facebook.github.io/react/)

### CSS

- [PostCSS](http://postcss.org/)
	- Compiler for CSS. Allows transforming CSS with JS. This gives you error messages for improperly formatted code and lets you specify plugins for css parsing like:
  - [CssNext](http://cssnext.io/)
    - Allows you to write CSS4 spec CSS. Allows variables, nested components. Similar to SASS/ LESS.
  	- [Autoprefixer](https://github.com/postcss/autoprefixer)
  		- Automatic vendor prefixing of CSS rules. No more having to worry about `-ms`, `-webkit` prefixes. Includes browser fixes for common broken css rules
- [Tachyons](http://tachyons.io/#getting-started)
  - Dead simple CSS framework for common classes.
