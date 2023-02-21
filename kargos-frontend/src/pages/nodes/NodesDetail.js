import React from 'react';
import { BrowserRouter as Router, Route, Redirect  } from 'react-router-dom';
import { NodesDetailPage } from '../../components/Nodes';

export default () => {
  return (
    <Router>
      <Route path="/nodes/detail/:page" component={page} />
    </Router>
  );
};

function page({match}) {
  const page = match.params;
  // Redirect default to page 1
  return (
    <NodesDetailPage page={page}></NodesDetailPage>
  );
}