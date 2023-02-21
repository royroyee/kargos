import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Nav, Pagination, Card, Table } from '@themesberg/react-bootstrap';
import { faGlobe, faInfoCircle, faExclamationCircle, faExclamationTriangle, faBug, faTerminal } from '@fortawesome/free-solid-svg-icons';

import React, { useState, useEffect } from 'react';


/**
 * Generate overall page for the Events.
 * @returns A JSX Component that shows events table, tab and pagenation.
 */
export const EventsOverviewPage = () => {
    return <EventsTableSection></EventsTableSection>
}

/**
 * This JSX will generate a table and table with pagenavition.
 * @returns A JSX that prints out the table with tabs.
 */
const EventsTableSection = () => {
    // Generate data to store in the later context.
    const [page, setPage] = useState(1);
    const [data, setData] = useState([]);
    const [select, setSelect] = useState('all');

    /**
     * This function will update table data using setData and refresh the table for the user.
     * This uses REST API /events/ from backend server.
     * @param {Number} argPage The page to query for.
     * @param {String} argSelection The selection of type to query.
     */
    function updateTable (argPage, argSelection) {
        var url;
        if (['Normal', 'Warning', 'Debug', 'System', 'Error'].includes(argSelection)) {
            url = "/events/?event=" + argSelection.toLowerCase() + "&page=" + argPage + "&per_page=10";
        } else {
            url = "/events/?&page=" + argPage + "&per_page=10";
        }

        var requestOptions = {
            method: 'GET',
            redirect: 'follow'
        };

        fetch(url, requestOptions)
            .then(response => response.text())
            .then(result => {
                setData(JSON.parse(result));
            })
            .catch(error => {
                console.log("empty response")
        });
    }

    /**
     * This function will handle selection click. Also triggers the code to re-render the table as well.
     * @param {String} selection The selection type to be used later.
     */
    function handleSelectionClick(selection) {
        setSelect(selection);
        setPage(1);
        updateTable(1, selection); // Need to send selection since the setSelect might not be updated yet.
    }
    
    /**
     * This function will handle new page click from pagination. Also will trigger the code to re-render the table.
     * @param {int} newPage The new page to set.
     */
    function handlePageinationClick(newPage) {
        setPage(newPage);
        updateTable(newPage, select); // Need to send selection since the setPage might not be updated yet.
    }

    /**
     * Generate Pageination for Events.
     * @todo add support for relative pagenation. For example, when we have 100 tabs, we need to have 25 ~ 35 printed out in the screen. But for now, it does not.
     * @returns JSX component that implements pageination.
     */
    const EventsPagination = () => {    
        const [count, setCount] = useState([]);
        const getCount = () => {
            var url;
            if (select == 'all') {
                url = "/events/count";
            } else {
                url = "/events/count?level=" + select.toLocaleLowerCase();
            }

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

        var eventsCount = count['count']

        const items = [];
        var totalPages = Math.ceil(eventsCount / 10) + 1;      
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
     * Generate a tab that offers users with selection of types of events to view.
     * @returns A JSX Component that performs tab action.
     */
    const NavBarSection = () => {
        return (
            <>
            <Nav fill defaultActiveKey="all" variant="pills" className="flex-column flex-md-row">
            <Nav.Item>
                <Nav.Link eventKey="all" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('all')} active={select === 'all'}>
                <FontAwesomeIcon icon={faGlobe} className="me-2" /> All
                </Nav.Link>
            </Nav.Item>
            <Nav.Item>
                <Nav.Link eventKey="normal" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('Normal')} active={select === 'Normal'}>
                <FontAwesomeIcon icon={faInfoCircle} className="me-2" /> Normal
                </Nav.Link>
            </Nav.Item>
            <Nav.Item>
                <Nav.Link eventKey="warning" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('Warning')} active={select === 'Warning'}>
                <FontAwesomeIcon icon={faExclamationCircle} className="me-2" /> Warning
                </Nav.Link>
            </Nav.Item>
            <Nav.Item>
                <Nav.Link eventKey="error" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('Error')} active={select === 'Error'}>
                <FontAwesomeIcon icon={faExclamationTriangle} className="me-2" /> Error
                </Nav.Link>
            </Nav.Item>
            <Nav.Item>
                <Nav.Link eventKey="debug" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('Debug')} active={select === 'Debug'}>
                <FontAwesomeIcon icon={faBug} className="me-2" /> Debug
                </Nav.Link>
            </Nav.Item>
            <Nav.Item>
                <Nav.Link eventKey="system" href="#" className="mb-sm-3 mb-md-0" onClick={() => handleSelectionClick('System')} active={select === 'System'}>
                <FontAwesomeIcon icon={faTerminal} className="me-2" /> System
                </Nav.Link>
            </Nav.Item>
            </Nav>
            </>
        );
    }

    /**
     * Generate table rows for each elements
     * @param {*} props The props that includes each data.
     * @returns A JSX Component that represents a set of Rows in table.
     */
    const TableRow = (props) => {  
        const { created, event_level, type, name, status, message } = props;
        const statusVariant = event_level === "Normal" ? "success"
            : event_level === "Warning" ? "warning"
            : event_level === "Debug" ? "info" 
            : event_level === "Error" ? "danger"
            : "primary";

        return (
            <>
            <tr>
            <td>
                <span className="fw-normal">
                {created}
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
                <span className={`fw-normal text-${statusVariant}`}>
                {event_level}
                </span>
            </td>
            <td>
                <span className="fw-normal">
                {status}
                </span>
            </td>
            <td>
                <span className="fw-normal">
                {message}
                </span>
            </td>
            </tr>
            </>
        );
    };

    /**
     * Generate a table for events.
     * @returns A JSX Component that represents events table.
     */
    const EventsTable = () => {
        return (
        <Table hover className="user-table align-items-center">
            <thead>
            <tr>
                <th style={{ width: '150px' }} className="border-bottom">TimeStamp</th>
                <th style={{ width: '150px' }} className="border-bottom">Type</th>
                <th style={{ width: '150px' }} className="border-bottom">Name</th>
                <th style={{ width: '150px' }} className="border-bottom">Level</th>
                <th style={{ width: '150px' }} className="border-bottom">Status</th>
                <th style={{ width: '150px' }} className="border-bottom">Message</th>
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
        );
    }

    return (
        <>
        <Card border="light" className="table-wrapper table-responsive shadow-sm">
            <Card.Header>
                <NavBarSection></NavBarSection>
            </Card.Header>
            <Card.Body className="pt-0">
                <EventsTable></EventsTable>
            </Card.Body>
            <Card.Footer>
                <EventsPagination></EventsPagination>
            </Card.Footer>
        </Card>
        </>
    );
}