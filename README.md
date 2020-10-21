# Traefik Datadog Event

This plugin can be used to generate Datadog events if certain patterns match.

## Plugin options

**APIKey**

*Required: true*

Your Datadog API key.

### Patterns

**Code**

This pattern compares the user defined status code, with the response code.

### Event options

**Title**

The event title. Limited to 100 characters. 

**Message**

The body of the event. Limited to 4000 characters. The text supports markdown.

**Priority**

*Default: normal*

The priority of the event. For example, normal or low. Allowed enum values: normal,low
