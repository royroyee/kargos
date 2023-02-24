import React from 'react';
import { BrowserRouter as Router, Route, Redirect  } from 'react-router-dom';
import { PodDetailPage } from '../../components/Pod';


export default () => {
  return (
    <Router>
      <Route path="/resources/pods/detail/:namespace/:page" component={page} />
    </Router>
  );
};

function page({match}) {
  const params = match.params;
  const name = params.page;
  const namespace = params.namespace;
  console.log(namespace + name)
  return <PodDetailPage page={name} namespace={namespace}></PodDetailPage>
}