package shared

import (
	"fmt"
	"sync"
)

var (
	oneErrorFactory sync.Once
	Error           ErrorFactory
)

var (
	errorTemplate          = `{"Status": "%d", "detail":"%s"}`
)

func init() {
	oneErrorFactory.Do(func() {
		Error = &errorFactory{}
	})
}

type ErrorFactory interface {
	InvalidPath(path, detail string) error
	InvalidFilter(filter, detail string) error
	InvalidType(path, expect, got string) error
	NoAttribute(path string) error
	MissingRequiredProperty(path string) error
	MutabilityViolation(path string) error
	InvalidParam(name, expect, got string) error
	InvalidParamGeneric() error
	ResourceNotFound(id, version string) error
	Duplicate(path string, value interface{}) error
	UnauthorisedRequest() error
	ForbiddenRequest() error
	DomainUnverified() error
	RequestCancelled() error
	Datastore(reason error) error
	Text(template string, args ...interface{}) error
}

type errorFactory struct{}

func (f *errorFactory) Text(template string, args ...interface{}) error {
	return fmt.Errorf(template, args...)
}

func (f *errorFactory) InvalidPath(path, detail string) error {
	return &InvalidPathError{path, detail}
}

// Invalid Path
type InvalidPathError struct {
	Path   string
	Detail string
}

func (e InvalidPathError) Error() string {
	return fmt.Sprintf("Path [%s] is invalid: %s", e.Path, e.Detail)
}

func (f *errorFactory) InvalidFilter(filter, detail string) error {
	return &InvalidFilterError{filter, detail}
}

// Invalid Filter
type InvalidFilterError struct {
	Filter string
	Detail string
}

func (e InvalidFilterError) Error() string {
	if len(e.Filter) > 0 {
		return fmt.Sprintf("Filter [%s] is invalid: %s", e.Filter, e.Detail)
	} else {
		return fmt.Sprintf("Filter is invalid: %s", e.Detail)
	}
}

func (f *errorFactory) InvalidType(path, expect, got string) error {
	return &InvalidTypeError{path, expect, got}
}

// Invalid Type
type InvalidTypeError struct {
	Path   string
	Expect string
	Got    string
}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("Invalid type at '%s', expected '%s', got '%s'", e.Path, e.Expect, e.Got)
}

func (f *errorFactory) NoAttribute(path string) error {
	return &NoAttributeError{path}
}

// No Attribute
type NoAttributeError struct {
	Path string
}

func (e *NoAttributeError) Error() string {
	return fmt.Sprintf("No attribute defined for path (segment) '%s'", e.Path)
}

func (f *errorFactory) MissingRequiredProperty(path string) error {
	return &MissingRequiredPropertyError{path}
}

// Missing Required Property

type MissingRequiredPropertyError struct {
	Path string
}

func (e *MissingRequiredPropertyError) Error() string {
	return fmt.Sprintf("Missing required property value at '%s'", e.Path)
}

func (f *errorFactory) MutabilityViolation(path string) error {
	return &MutabilityViolationError{path}
}

// Mutability Violation Error

type MutabilityViolationError struct {
	Path string
}

func (e *MutabilityViolationError) Error() string {
	return fmt.Sprintf("Violated mutability rule at '%s'", e.Path)
}

func (f *errorFactory) InvalidParam(name, expect, got string) error {
	return &InvalidParamError{name, expect, got}
}

// Invalid Param Error

type InvalidParamError struct {
	Name   string
	Expect string
	Got    string
}

func (e *InvalidParamError) Error() string {
	return fmt.Sprintf("Invalid parameter for %s, expect %s, but got %s", e.Name, e.Expect, e.Got)
}

func (f *errorFactory) InvalidParamGeneric() error {
	return &InvalidParamGenericError{}
}

// Invalid Param Error Generic
type InvalidParamGenericError struct {
}

func (e *InvalidParamGenericError) Error() string {
	return fmt.Sprintf("Invalid parameter provided")
}

func (f *errorFactory) ResourceNotFound(id, version string) error {
	return &ResourceNotFoundError{id, version}
}

// Resource Not Found

type ResourceNotFoundError struct {
	Id      string
	Version string
}

func (e ResourceNotFoundError) Error() string {
	if len(e.Id) > 0 {
		if len(e.Version) > 0 {
			return fmt.Sprintf("Resource not found for id '%s' and version '%s'", e.Id, e.Version)
		}
		return fmt.Sprintf("Resource not found for id '%s'", e.Id)
	}
	return "Resource not found"
}

func (f *errorFactory) Duplicate(path string, value interface{}) error {
	return &DuplicateError{path, value}
}

// Duplicate Error
type DuplicateError struct {
	Path  string
	Value interface{}
}

func (e DuplicateError) Error() string {
	return fmt.Sprintf("Resource has duplicate value '%v' at path '%s'", e.Value, e.Path)
}

// Unauthorised Error
type UnauthorisedError struct {
}

func (e UnauthorisedError) Error() string {
	return "Unauthorised: Username or password improper"
}

func (f *errorFactory) UnauthorisedRequest() error {
	return &UnauthorisedError{}
}

type ForbiddenError struct {
}

func (e ForbiddenError) Error() string {
	return "Forbidden: Token invalid"
}

func (f *errorFactory) ForbiddenRequest() error {
	return &ForbiddenError{}
}

type UnverifiedDomain struct {
}

func (e UnverifiedDomain) Error() string {
	return "Domain is unverified"
}

func (f *errorFactory) DomainUnverified() error {
	return &UnverifiedDomain{}
}

type RequestCancelled struct {
}

func (e RequestCancelled) Error() string {
	return "Request was cancelled by the callee"
}

func (f *errorFactory) RequestCancelled() error {
	return &RequestCancelled{}
}

type PaymentInvalidError struct {
	Reason string
}

func (e *PaymentInvalidError) Error() string {
	return fmt.Sprintf("No attribute defined for path (segment) '%s'", e.Reason)
}
func (f *errorFactory) PaymentInvalid(reason string) error {
	return &PaymentInvalidError{reason}
}

type DatastoreError struct {
	Reason error
}

func (e *DatastoreError) Error() string {
	return fmt.Sprintf("Error in database reason %s", e.Reason)
}

func (f *errorFactory) Datastore(reason error) error {
	return &DatastoreError{reason}
}
