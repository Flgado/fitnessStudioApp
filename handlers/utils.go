package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
// @param urlValues url.Values - URL query parameters.
//
// @return api.ClasseFilters - Populated struct containing class filters.
// @return error - Error if any parameter fails to parse.
func buildClassFilters(urlValues url.Values) (api.ClasseFilters, error) {
	filters := api.ClasseFilters{
		Name: urlValues.Get("class_name"),
	}

	// Parse start date greater than or equal to
	if startDateStr := urlValues.Get("start_date"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "start_date")
		}
		filters.StartDateGte = &startDate
	}

	// Parse end date less than or equal to
	if endDateStr := urlValues.Get("end_date"); endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "end_date")
		}
		filters.EndDateLe = &endDate
	}

	// Parse capacity greater than or equal to
	if capacityGteStr := urlValues.Get("capacity_gte"); capacityGteStr != "" {
		capacityGte, err := strconv.Atoi(capacityGteStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "capacity_gte")
		}
		filters.CapacityGte = &capacityGte
	}

	// Parse capacity less than or equal to
	if capacityLeStr := urlValues.Get("capacity_le"); capacityLeStr != "" {
		capacityLe, err := strconv.Atoi(capacityLeStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "capacity_le")
		}
		filters.CapacityLe = &capacityLe
	}

	// Parse number of registrations greater than or equal to
	if numRegistrationsGteStr := urlValues.Get("num_registrations_gte"); numRegistrationsGteStr != "" {
		numRegistrationsGte, err := strconv.Atoi(numRegistrationsGteStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "num_registrations_gte")
		}
		filters.NumRegistrationsGte = &numRegistrationsGte
	}

	// Parse number of registrations less than or equal to
	if numRegistrationsLeStr := urlValues.Get("num_registrations_le"); numRegistrationsLeStr != "" {
		numRegistrationsLe, err := strconv.Atoi(numRegistrationsLeStr)
		if err != nil {
			return api.ClasseFilters{}, buildFormatParameterError(err, "num_registrations_le")
		}
		filters.NumRegistrationsLe = &numRegistrationsLe
	}

	return filters, nil
}

// buildFormatParameterError constructs a formatted error for wrong parameter format.
//
// This function takes an error describing the failure to parse a parameter,
// the name of the parameter, and constructs a formatted error response.
// It returns a utils.Error struct with appropriate status code, message, and details.
//
// @param e error - Error describing the failure to parse a parameter.
// @param p string - Name of the parameter that failed to parse.
//
// @return utils.Error - Formatted error response.
func buildFormatParameterError(e error, p string) utils.Error {
	return utils.E(http.StatusBadRequest,
		e,
		map[string]string{"message": "Wrong parameter pass"},
		fmt.Sprintf("Wrong parameter pass as %s", p),
		"Use the right parameter value. Read documentation for more details")
}

func buildError(e error, m string, s string) utils.Error {
	return utils.E(http.StatusBadRequest,
		e,
		map[string]string{"message": m},
		"",
		s)
}
