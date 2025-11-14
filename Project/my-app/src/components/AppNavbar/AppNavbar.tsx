// src/components/AppNavbar/AppNavbar.tsx
import { Navbar, Nav, Container, NavDropdown } from "react-bootstrap";
import { LinkContainer } from "react-router-bootstrap"; // чтобы использовать <Nav.Link> с маршрутом

export const AppNavbar = () => {
  return (
    <Navbar bg="light" expand="lg">
      <Container>
        <LinkContainer to="/">
          <Navbar.Brand>My Music App</Navbar.Brand>
        </LinkContainer>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="me-auto">
            <LinkContainer to="/">
              <Nav.Link>Home</Nav.Link>
            </LinkContainer>
            <LinkContainer to="/itunes">
              <Nav.Link>iTunes</Nav.Link>
            </LinkContainer>
            <NavDropdown title="Albums" id="basic-nav-dropdown">
              <LinkContainer to="/album/1">
                <NavDropdown.Item>Album 1</NavDropdown.Item>
              </LinkContainer>
              <LinkContainer to="/album/2">
                <NavDropdown.Item>Album 2</NavDropdown.Item>
              </LinkContainer>
              <NavDropdown.Divider />
              <LinkContainer to="/album/3">
                <NavDropdown.Item>Album 3</NavDropdown.Item>
              </LinkContainer>
            </NavDropdown>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};
