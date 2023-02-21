
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Row, Col, Nav, Form, Pagination, Dropdown, ButtonGroup, Card, Button, Table, Alert } from '@themesberg/react-bootstrap';
import { faEllipsisH, faChevronRight, faTools, faSmile } from '@fortawesome/free-solid-svg-icons';

import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { RightSlidePage } from './RightSlide';
import AccordionComponent from '../components/AccordionComponent'

import deploymentInfo from "../fakeData/DUMMY_DEPLOYMENT_INFO"

/**
 * Generates overview page for workload.
 * @param {*} props The props information.
 * @returns The overview page for workload.
 */
export const WorkloadOverviewPage = () => {
  // Generate data to store in the later context.
  const [page, setPage] = useState(1);
  const [data, setData] = useState([]);
  const [select, setSelect] = useState('all');
  const [namespace, setNamespace] = useState('default');

  /**
   * Handles option change when user clicked form selection component.
   * @param {event} event The event that user clicked in forms.
   */
  function handleOptionChange(event) {
    const namespaceSelection = event.target.value;
    setNamespace(namespaceSelection);
    updateTable(namespaceSelection, 1, 'all');
  }

  /**
   * Generates Namespace selection form.
   * @returns JSX Component that has form attribute that can make user chose which namespace to view.
   */
  const NamespaceSelection = () => {
    const [namespaces, setNamespaces] = useState([]);

    // Retrieve data from REST API.
    const getNamespaces = () => {
      const url = "/workload/namespaces";
      var requestOptions = {
          method: 'GET',
          redirect: 'follow'
      };
  
      fetch(url, requestOptions)
          .then(response => response.text())
          .then(result => {
            setNamespaces(JSON.parse(result));
          })
          .catch(error => console.log('error', error));
  }
    useEffect(() => {
      getNamespaces();
    }, []);
      
    // Generate options for forms.
    var options = [];
    for (var i = 0; i < namespaces.length; i++) {
      options.push(<option>{namespaces[i]}</option>);
    }

    // Return namespace selection.
    return (
      <>
        <Form>
          <Form.Group className="mb-3">
            <Form.Label>Namespace</Form.Label>
            <Form.Select value={namespace} onChange={handleOptionChange}>
              {options}
            </Form.Select>
          </Form.Group>
          </Form>
      </>
    );
  }

  /**
   * A function that handles user's selection.
   * @param {String} userSelection The user's selection
   */
  function handleSelectionClick(userSelection) {
    setSelect(userSelection);
    updateTable(namespace, 1, userSelection);
  }
  
  /**
     * Generate a tab that offers users with selection of types of events to view.
     * @returns A JSX Component that performs tab action.
     */
  const NavBarSection = () => {
    return (
        <>
        <Nav fill defaultActiveKey="all" variant="pills" className="flex-column flex-md-row">
        <Nav.Item>
            <Nav.Link eventKey="all" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('all')} active={select === 'all'}>
              All
            </Nav.Link>
        </Nav.Item>
        <Nav.Item>
            <Nav.Link eventKey="replicaset" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('ReplicaSet')} active={select === 'ReplicaSet'}>
              ReplicaSet
            </Nav.Link>
        </Nav.Item>
        <Nav.Item>
            <Nav.Link eventKey="deployment" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('Deployment')} active={select === 'Deployment'}>
             Deployment
            </Nav.Link>
        </Nav.Item>
        <Nav.Item>
            <Nav.Link eventKey="daemonset" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('DaemonSet')} active={select === 'DaemonSet'}>
              DaemonSet
            </Nav.Link>
        </Nav.Item>
        <Nav.Item>
            <Nav.Link eventKey="statefulset" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('StatefulSet')} active={select === 'StatefulSet'}>
              StatefulSet
            </Nav.Link>
        </Nav.Item>
        <Nav.Item>
            <Nav.Link eventKey="job" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('Job')} active={select === 'Job'}>
              Job
            </Nav.Link>
        </Nav.Item>
        <Nav.Item>
            <Nav.Link eventKey="cronjob" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('CronJob')} active={select === 'CronJob'}>
              CronJob
            </Nav.Link>
        </Nav.Item>
        </Nav>
        </>
    );
  }

  /**
   * This function will handle new page click from pagination. Also will trigger the code to re-render the table.
   * @param {int} newPage The new page to set.
   */
  function handlePageinationClick(newPage) {
    setPage(newPage);
    updateTable(namespace, newPage, select);
  }

  /**
     * Generate Pageination for Workload.
     * @todo add support for relative pagenation. For example, when we have 100 tabs, we need to have 25 ~ 35 printed out in the screen. But for now, it does not.
     * @returns JSX component that implements pageination.
     */
  const WorkloadPagination = () => {
    const [count, setCount] = useState([]);
    const getCount = () => {
        var url; 
        if (select == 'all') {
          url = "/workload/count?namespace=" + namespace;
        } else {
          url = "/workload/count?namespace=" + namespace + "&type=" + select.toLowerCase();
        }

        console.log(url);

        var requestOptions = {
            method: 'GET',
            redirect: 'follow'
        };
    
        fetch(url, requestOptions)
            .then(response => response.text())
            .then(result => {
                setCount(JSON.parse(result));
            })
            .catch(error => console.log('error', error));
    }
        
    useEffect(() => {
        getCount();
    }, []);

    
    var workloadCount = count['count']    

    console.log("TOTAL COUNT : " + workloadCount)
    const items = [];
    var totalPages = Math.ceil(workloadCount / 10) + 1;      
    var prevDisabled = (page == 1);
    var nextDisabled = (page == totalPages - 1) || totalPages == 1;

    items.push(
        <Pagination.Prev disabled={prevDisabled} onClick={() => handlePageinationClick((page - 1))}>
        Previous
        </Pagination.Prev>
    );
    
    for (let i = 1; i < totalPages; i++) {
        if (i != Number(page)) {
        items.push(<Pagination.Item onClick={() => handlePageinationClick(i)}>
            {i}
        </Pagination.Item>);
        } else {
        items.push(<Pagination.Item active>
            {i}
        </Pagination.Item>);
        }
    }
    
    items.push(
        <Pagination.Next disabled={nextDisabled} onClick={() => handlePageinationClick((page + 1))}>
        Next
        </Pagination.Next>
    );
    
    return (
        <Pagination className="mb-2 mb-lg-0">
            {items}
        </Pagination>
    );
  }

  /**
   * This function will update table from data using setData and refresh the table for the user.
   * This uses REST API /workload/ API from backend.
   */
  function updateTable (argNamespace, argPage, argSelect) {
    var url;
    if (argSelect == 'all') {
      url = "/workload/?namespace=" + argNamespace + "&page=" + argPage + "&per_page=10"
    } else {
      url = "/workload/?namespace=" + argNamespace + "&controller=" + argSelect.toLowerCase() + "&page=" + argPage + "&per_page=10"
    }

    console.log(url)
    var requestOptions = {
        method: 'GET',
        redirect: 'follow'
    };

    fetch(url, requestOptions)
        .then(response => response.text())
        .then(result => {
            setData(JSON.parse(result));
        })
        .catch(error => console.log('error', error));
}

useEffect(() => {
    updateTable(namespace, page, select);
}, []);
  
  return (
    <article>    
      <Card border="light" className="table-wrapper table-responsive shadow-sm">
        <Card.Header>
        <NamespaceSelection></NamespaceSelection>
          <NavBarSection></NavBarSection>
        </Card.Header>
        <Card.Body className="pt-0">
          <WorkloadOverviewTable data={data}></WorkloadOverviewTable>
        </Card.Body>
        <Card.Footer>
          <WorkloadPagination></WorkloadPagination>
        </Card.Footer>
        </Card>
    </article>
  );
}


/**
 * Generate a information page for a workload. This is meant to be displayed in the half of the page.
 * @todo Add REST API support that retrieves data from backend.
 * @param {props} props The props to use.
 * @returns A JSX Component that represents information for a single workload.
 */
const WorkloadDetail = (props) => {
  const { name, type, namespace } = props;

  /**
   * Generate table rows for each elements
   * @param {*} props The props that includes each data.
   * @returns A JSX Component that represents a set of Rows in table.
  */
  const TableRow = (props) => {
    const { type, value } = props;
    return (
      <tr>
        <td>
          <span className="fw-normal">
            {type}
          </span>
        </td>
        <td>
          <span className="fw-normal">
            {value}
          </span>
        </td>
      </tr>
    );
  };

  /**
   * Generate Table that contains controller information.
   * @returns A JSX Component that represents the table for information on controller.
   */
  const WorkloadInfoTable = () => {
    return (
      <>
        <Table hover className="user-table align-items-center">
          <thead>
            <div className="d-block">
              <h5 className="fw-normal mb-2">
                Workload Information
              </h5>
            </div>
            <tr>
              <th className="border-bottom">Type</th>
              <th className="border-bottom">Value</th>
            </tr>
          </thead>
          <tbody>
            {deploymentInfo.controllerInfo.map(t => <TableRow key={`pods-${t.index}`} {...t} />)}
          </tbody>
        </Table>
      </>
    );
  }

  /**
   * Generate a dropdown with given data.
   * @returns A JSX Component that represents dropodown .
   */
  const DropDownSection = (props) => {
    const { data, name } = props;
    var accordionElements = [];
    for (var i = 0 ; i < data.length ; i++) {
      const name = data[i].name;
      const containers = data[i].info;

      accordionElements.push({
        id: name,
        eventKey: name,
        title: name,
        description: 
        <Card border="light" className="table-wrapper table-responsive shadow-sm">
          <Table hover className="user-table align-items-center">
            <thead>
              <tr>
                <th className="border-bottom">Type</th>
                <th className="border-bottom">Value</th>
              </tr>
            </thead>
            <tbody>
              {containers.map(t => <TableRow key={`processes-${t.index}`} {...t} />)}
            </tbody>
          </Table>
      </Card>
      });
    }

    return (
      <>
        <div className="d-block">
          <h5 className="fw-normal mb-2">
            {name}
          </h5>
        </div>
        <Card border="light" className="table-wrapper table-responsive shadow-sm">
            <AccordionComponent
              defaultKey = "none"
              data = {
                accordionElements
              }
            />
        </Card>
      </>
    );
  }

  /**
   * Generate a conditions table with given data.
   * @returns A JSX Component that represents conditions table.
   */
  const ConditionsSection = (props) => {
    const { data } = props;

    /**
     * Generate table rows for each elements
     * @param {*} props The props that includes each data.
     * @returns A JSX Component that represents a set of Rows in table.
    */
    const TableRow = (props) => {
      const { type, status, reason } = props;
      return (
        <tr>
          <td>
            <span className="fw-normal">
              {type}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {status}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {reason}
            </span>
          </td>
        </tr>
      );
    };

      return (
        <>
          <div className="d-block">
            <h5 className="fw-normal mb-2">
              Conditions
            </h5>
          </div>
          <Table hover className="user-table align-items-center">
            <thead>
            <tr>
                <th className="border-bottom">Type</th>
                <th className="border-bottom">Status</th>
                <th className="border-bottom">Reason</th>
              </tr>
            </thead>
            <tbody>
              {data.map(t => <TableRow {...t} />)}
            </tbody>
          </Table>
        </>
      );
  }

  /**
  // Return the whole page.
  return (
    <>
      <Row className="justify-content-md-center">
        <div className="d-block">
          <h3 className="fw-normal mb-2">
            {name}
          </h3>
        </div>
        <br></br>
          <Col className="mb-4">
            <Card border="light" className="table-wrapper table-responsive shadow-sm">
              <Card.Body className="pt-0">
                <WorkloadInfoTable></WorkloadInfoTable>
                <br></br>
                <DropDownSection data name="Template Containers"></DropDownSection>
                <br></br>
                <DropDownSection data name="Volumes"></DropDownSection>
                <br></br>
                <ConditionsSection data={deploymentInfo.conditions}></ConditionsSection>
              </Card.Body>
            </Card>
          </Col>
      </Row>
    </>
  );

  **/

  return (
    <>
      <Row className="justify-content-md-center">
        <Col className="mb-4">
          <div className="d-block">
            <h3 className="fw-normal mb-2">
              {name}
            </h3>
          </div>
          <TempWorkloadInfo name={name} namespace={namespace} type={type}></TempWorkloadInfo>
          <TempVolumeContainerDropDown name={name} namespace={namespace}></TempVolumeContainerDropDown>
          </Col>
      </Row>
    </>
  );
};

const TempWorkloadInfo = (props) => {
  const { name, namespace, type } = props;
  const [ data, setData ] = useState({ "template_containers":[], "volumes":[] });

  // Retrieve data from REST API.
  // This will retrieve volume and template container information.
  const getData = () => {
    const url = "/workload/info/" + namespace + '/' + name + "?type=" + type;
    var requestOptions = {
        method: 'GET',
        redirect: 'follow'
    };

    fetch(url, requestOptions)
        .then(response => response.text())
        .then(result => {
          setData(JSON.parse(result));
        })
        .catch(error => console.log('error', error));
  }

  useEffect(() => {
    getData();
  }, []);

  const TableRow = (props) => {
    const { type, value } = props;
    // TODO ADD LINK LIKE CARD TO THE NODE's NAME and delete action
    return (
      <tr>
        <td>
          <span className="fw-normal">
            {type}
          </span>
        </td>
        <td>
          <span className="fw-normal">
            { Array.isArray(value) ? value.join(", ") : value}
          </span>
        </td>
      </tr>
    );
  };
  
  // Generate Table rows by iterating key and values.
  var items = [];
  const keys = Object.keys(data);

  // Iterate over keys and generate table rows.
  for (const key of keys) {
      items.push(<TableRow type={key} value={data[key]} />);
  }

  return (
    <>
      <div className="d-block">
        <h5 className="fw-normal mb-2">
          Workload Information
        </h5>
      </div>
      <Table hover className="user-table align-items-center">
        <thead>
          <tr>
            <th className="border-bottom">Type</th>
            <th className="border-bottom">Value</th>
          </tr>
        </thead>
        <tbody>
          {items}
        </tbody>
      </Table>
    </>
  );
}

/**
 * Generate a temporary volume and container dropdown. This is for temp usage
 * @todo Remove this function to automatically show volume and container dropdown using <DropDownSection>
 * @deprecated
 * @param {props} props The props.
 * @returns Container and volume dropdown. 
 */
const TempVolumeContainerDropDown = (props) => {
  const { name, namespace } = props;
  const [ data, setData ] = useState({ "template_containers":[], "volumes":[] });

  // Retrieve data from REST API.
  // This will retrieve volume and template container information.
  const getData = () => {
    const url = "/workload/detail/" + namespace + '/' + name;
    var requestOptions = {
        method: 'GET',
        redirect: 'follow'
    };

    fetch(url, requestOptions)
        .then(response => response.text())
        .then(result => {
          setData(JSON.parse(result));
        })
        .catch(error => console.log('error', error));
  }

  useEffect(() => {
    getData();
  }, []);

  /**
   * Generate a dropdown with given data. This is for tmp usage only
   * @todo Delete this and use <DropDownSection> instead
   * @deprecated
   * @returns A JSX Component that represents dropodown .
   */
  const DropDownSection = (props) => {
    const { data, name } = props;
    var accordionElements = [];
    for (var i = 0 ; i < data.length ; i++) {
      accordionElements.push({
        id: data[i],
        eventKey: data[i],
        title: data[i],
        description: 
          <center>
            <Alert variant="info">
              <FontAwesomeIcon icon={faTools} className="icon-dark" />  We are working on this... Please hang tight! <FontAwesomeIcon icon={faSmile} className="icon-dark" />
            </Alert>
          </center>
      });
    }

    return (
      <>
        <div className="d-block">
          <h5 className="fw-normal mb-2">
            {name}
          </h5>
        </div>
          <AccordionComponent
            defaultKey = "none"
            data = {
              accordionElements
            }
          />
      </>
    );
  }

  return (
      <>
        <Row className="justify-content-md-center">
          <br></br>
            <Col className="mb-4">
              <br></br>
              <DropDownSection data={data.template_containers} name="Template Containers"></DropDownSection>
              <br></br>
              <DropDownSection data={data.volumes} name="Volumes"></DropDownSection>
              <br></br>
            </Col>
        </Row>
      </>
    );
}

/**
 * Generates a table of overviewed 
 * @param {*} props The props information.
 * @returns The worload overview table.
 */
const WorkloadOverviewTable = (props) => {
    const { data } = props; 
    const TableRow = (props) => {
      const { namespace, type, name, pods } = props;
      function generateDropDown(pod) {
        return (
          <Dropdown.Item>
            <Link to={"/resources/pods/detail/" + namespace + "/" + pod}> {pod} </Link>
          </Dropdown.Item>
        );
      };

      // TODO ADD LINK LIKE CARD TO THE NODE's NAME and delete action
      return (
        <tr>
          <td>
            <span className="fw-normal">
              {namespace}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {type}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              {name}
            </span>
          </td>
          <td>
            <span className="fw-normal">
              <Dropdown as={ButtonGroup}>
                <Dropdown.Toggle as={Button} split variant="link" className="text-dark m-0 p-0">
                  <span className="icon icon-sm">
                    <FontAwesomeIcon icon={faEllipsisH} className="icon-dark" />
                  </span>
                </Dropdown.Toggle>
                <Dropdown.Menu>
                  {pods.map(pod => (generateDropDown(pod)))}
                </Dropdown.Menu>
              </Dropdown>
            </span>
          </td>
          <td>
            <span className="fw-normal">
              <RightSlidePage clickButton={<FontAwesomeIcon icon={faChevronRight}></FontAwesomeIcon>} content={<WorkloadDetail name={name} type={type} namespace={namespace}></WorkloadDetail>}> </RightSlidePage>
            </span>
          </td>
        </tr>
      );
    };
  
    return (
      <>
      <Card border="light" className="table-wrapper table-responsive shadow-sm">
          <Card.Body className="pt-0">
            <Table className="user-table align-items-center">
              <thead>
                <tr>
                  <th className="border-bottom">Namespace</th>
                  <th className="border-bottom">Type</th>
                  <th className="border-bottom">Name</th>
                  <th className="border-bottom">Pods</th>
                  <th className="border-bottom">Detail</th>
                </tr>
              </thead>
              <tbody>
              { // If the retrieved data was null, set it empty table row.
                data ? (
                <>
                    {data.map(t => (<TableRow {...t} />
                    ))} 
                </>
                ) : (
                    <tr>
                    </tr>
                )
            }
              </tbody>
            </Table>
          </Card.Body>
        </Card>
      </>
    );
  }