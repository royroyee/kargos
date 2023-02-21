import React from 'react';
import { Col, Row } from '@themesberg/react-bootstrap';
import { EventsOverviewPage } from '../../components/Events';

export default () => {
  return (
    <article>
          <Row className="d-flex flex-wrap flex-md-nowrap py-4">
            <Col className="d-block mb-4 mb-md-0">
              <h1 className="h2">Events</h1>
              <p className="mb-0">
                Monitor Kubernetes cluster events.
              </p>
              </Col>
          </Row>
       <EventsOverviewPage></EventsOverviewPage>
    </article>
  );
};