gorets explorer client
======================

## Dependencies

Requires node >= 6.3.1. Can switch by running `nvm use`.

## Installing

`npm i`

## Running

`npm run start`
Visit http://localhost:8000

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
