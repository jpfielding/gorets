import React from 'react';
import ReactDOM from 'react-dom';
import 'styles/app.css';

ReactDOM.render(
  <div className="helvetica">
    <h1 className="f1">gorets react</h1>
    <div className="flex flex-row justify-around">
      <div className="box tc">test</div>
      <div className="box tc">test</div>
      <div className="box tc">test</div>
    </div>
  </div>,
  document.getElementById('app')
);
