package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	api "github.com/Flgado/fitnessStudioApp/internal/api/models"
	"github.com/Flgado/fitnessStudioApp/utils"
)

// buildClassFilters parses the URL query parameters to build class filters.
//
// This function takes a url.Values object containing query parameters and constructs
// an api.ClasseFilters struct with optional filtering parameters for classes.
// It returns a populated api.ClasseFilters struct and nil error if successful.
// If any parameter fails to parse, it returns an empty api.ClasseFilters struct
// and an error describing the issue.
//
// param urlValues url.Values - URL query parameters.
//
// return api.ClasseFilters - Populated struct containing class filters.
// return error - Error if any parameter fails to parse.
func buildClassFilters(urlValues url.Values) (api.ClasseFilters, error) {

	className := strings.TrimSpace(urlValues.Get("className"))
	filters := api.ClasseFilters{
		Name: className,
	}

	layout := "2006-01-02"
	// Parse start date greater than or equal to
	if startDateStr := urlValues.Get("startDate"); startDateStr != "" {
		startDate, err := time.Parse(layout, startDateStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "startDate")
		}
		filters.StartDateGte = &startDate
	}

	// Parse end date less than or equal to
	if endDateStr := urlValues.Get("endDate"); endDateStr != "" {
		endDate, err := time.Parse(layout, endDateStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "endDate")
		}
		filters.EndDateLe = &endDate
	}

	// Parse capacity greater than or equal to
	if capacityGteStr := urlValues.Get("capacityGte"); capacityGteStr != "" {
		capacityGte, err := strconv.Atoi(capacityGteStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "capacityGte")
		}
		filters.CapacityGte = &capacityGte
	}

	// Parse capacity less than or equal to
	if capacityLeStr := urlValues.Get("capacityLe"); capacityLeStr != "" {
		capacityLe, err := strconv.Atoi(capacityLeStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "capacityLe")
		}
		filters.CapacityLe = &capacityLe
	}

	// Parse number of registrations greater than or equal to
	if numRegistrationsGteStr := urlValues.Get("numRegistrationsGte"); numRegistrationsGteStr != "" {
		numRegistrationsGte, err := strconv.Atoi(numRegistrationsGteStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "numRegistrationsGte")
		}
		filters.NumRegistrationsGte = &numRegistrationsGte
	}

	// Parse number of registrations less than or equal to
	if numRegistrationsLeStr := urlValues.Get("numRegistrationsLe"); numRegistrationsLeStr != "" {
		numRegistrationsLe, err := strconv.Atoi(numRegistrationsLeStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "numRegistrationsLe")
		}
		filters.NumRegistrationsLe = &numRegistrationsLe
	}

	return filters, nil
}

func BuildUpdateClass(date *string, name *string, capacity *int) (api.UpdateClass, error) {
	layout := "2006-01-02"
	if date != nil {
		newDate, err := time.Parse(layout, *date)
		if err != nil {
			return api.UpdateClass{}, utils.E(http.StatusBadRequest,
				err,
				map[string]string{"message": "BadRequest"},
				"Request body not expected",
				"Read our documentation for more details")
		}

		return api.UpdateClass{
			Name:     name,
			Date:     &newDate,
			Capacity: capacity,
		}, nil
	}

	return api.UpdateClass{
		Name:     name,
		Date:     nil,
		Capacity: capacity,
	}, nil
}

// buildFormatParameterError constructs a formatted error for wrong parameter format.
//
// This function takes an error describing the failure to parse a parameter,
// the name of the parameter, and constructs a formatted error response.
// It returns a utils.Error struct with appropriate status code, message, and details.
//
// param e error - Error describing the failure to parse a parameter.
// param p string - Name of the parameter that failed to parse.
//
// return utils.Error - Formatted error response.
func buildFormatParameterError(e error, p string) utils.Error {
	return utils.E(http.StatusBadRequest,
		e,
		map[string]string{"message": "Wrong parameter pass"},
		fmt.Sprintf("Wrong parameter pass as %s", p),
		"Use the right parameter value. Read documentation for more details")
}
