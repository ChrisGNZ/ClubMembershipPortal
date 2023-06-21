insert into [Roles](RoleName, RoleDescription)
values ('Developer', 'System Developer - superuser access for development, support, and troubleshooting.')
insert into [Roles](RoleName, RoleDescription)
values ('Membership Manager',
        'Full Read/Write access to the Membership database - typically the Club Membership officer will have this role.')
insert into [Roles](RoleName, RoleDescription)
values ('Committee',
        'Read-nly access to the Membership database - committee members can read but not update the membership database')


insert into [WebFormQuestions]([FormID], [InputFieldName], [EntityName], [QuestionText], [QuestionType], [TemplateName],
                               [Seq], [SubSeq], [AnswerRequired], [TextRows], [RadioOption], [QuestionLabel],
                               [QuestionExtraText])
select 1,
       'FirstName',
       'FirstName',
       'First Name (Known As)',
       'text',
       'StdText',
       10,
       NULL,
       'Y',
       NULL,
       NULL,
       NULL,
       NULL
union
select 1,
       'LastName',
       'LastName',
       'Surname',
       'text',
       'StdText',
       20,
       NULL,
       'Y',
       NULL,
       NULL,
       NULL,
       NULL
union
select 1,
       'Address',
       'Address',
       'Address',
       'textarea',
       'StdTextArea',
       30,
       NULL,
       'Y',
       4,
       NULL,
       NULL,
       NULL
union
select 1,
       'Postcode',
       'Postcode',
       'Postcode',
       'number',
       'StdText',
       40,
       NULL,
       'N',
       NULL,
       NULL,
       NULL,
       NULL
union
select 1,
       'PreferredPhone',
       'PreferredPhone',
       'Preferred Phone or Mobile',
       'tel',
       'StdText',
       50,
       NULL,
       'Y',
       NULL,
       NULL,
       NULL,
       NULL
union
select 1,
       'Email',
       'Email',
       'Email',
       'email',
       'StdText',
       70,
       NULL,
       'Y',
       NULL,
       NULL,
       NULL,
       NULL
union
select 1,
       'EmergencyContact',
       'EmergencyContact',
       'Name and phone of Spouse/Partner',
       'text',
       'StdText',
       80,
       NULL,
       'N',
       NULL,
       NULL,
       'Name and phone of Spouse/Partner',
       'For emergency contact and so we can be polite when calling'
union
select 1,
       'Occupation',
       'Occupation',
       'Occupation',
       'text',
       'StdText',
       90,
       NULL,
       'N',
       NULL,
       NULL,
       'Occupation',
       'Required by Incorporated Societies Act - If you are retired, we would like to know the skills and occupation you had when working - enter: Retired (previous occupation)'
union
select 1,
       'NAWMembershipNumber',
       'NAWMembershipNumber',
       'NAW Membership #',
       'text',
       'StdText',
       120,
       NULL,
       'N',
       NULL,
       NULL,
       'NAW Membership #',
       '(leave blank if you are not a member of NAW)'
union
select 1,
       'AllowSharingOfMembershipDetails',
       'AllowSharingOfMembershipDetails',
       'Membership Listing?',
       'radio',
       'RadioButton',
       130,
       10,
       'Y',
       NULL,
       'YES',
       '',
       'We would like to have a member list for sharing among our members. The Privacy Act 2020 requires that you agree to this selecting YES. The member list will be in printed format only and given only to our club members and will contain only names, phone numbers and email addresses.'
union
select 1,
       'AllowSharingOfMembershipDetails',
       'AllowSharingOfMembershipDetails',
       'Membership Listing?',
       'radio',
       'RadioButton',
       140,
       20,
       'Y',
       NULL,
       'No',
       NULL,
       NULL
union
select 1,
       'YouthMember',
       'YouthMember',
       'Youth?',
       'radio',
       'RadioButton',
       150,
       10,
       'Y',
       NULL,
       'Y',
       'Youth?',
       'Subject to terms and conditions'
union
select 1,
       'YouthMember',
       'YouthMember',
       'Youth?',
       'radio',
       'RadioButton',
       160,
       20,
       'Y',
       NULL,
       'N',
       NULL,
       NULL
union
select 1, 'SubmittedFee', 'SubmittedFee', 'Submitted Fee', 'text', 'StdText', 'Y'
order by 7, 8

---------------------------------------------------------------------------------------------------------------------------------------------------------------
insert into WebFormQuestionTemplates([QuestionTemplateName], [QuestionTemplateText])
values ('RadioButtonDetail', '<div class="form-check form-check-inline mt-4">
    <input
            class="form-check-input"
            type="radio"
            name="{{ .radioButtonName }}"
            id="{{ .radioButtonID }}"
            value="{{ .radioButtonValue }}"
    />
    <label class="form-check-label" for="workstatusOption3">{{ .radioButtonValue }}</label>
</div>')

insert into WebFormQuestionTemplates([QuestionTemplateName], [QuestionTemplateText])
values ('RadioButtonHeader',
        '<div class="container-fluid">
            <div class="row">
                <div class="col-md-12 mb-4 border border-1 rounded">
                    <div class="row"><div class="col-md-12 mt-1"><span style="font-size: smaller;">{{ .InputFieldName }}</span></div></div>
                    <div class="row"><div class="col-md-12 mt-1"><span class="fw-lighter" style="font-size: smaller;">{{ .questionExtraText }}</span></div></div>')

insert into WebFormQuestionTemplates([QuestionTemplateName], [QuestionTemplateText])
values ('StdText', '<div class="row">
    <div class="col-md-12 mb-2">
        <div class="form-floating">
            <input type="{{ .questionType }}" name="{{ .InputFieldName }}" id="{{ .questionID }}" class="form-control" value="{{ .questionValue }}" {{ .required }} />
            <label class="form-label" for="{{ .questionID }}">{{ .questionLabel }}</label>
        </div>
    </div>
</div>
<div class="row"><div class="col-md-12 mb-4"><span class="fw-lighter" style="font-size: smaller;">{{ .questionExtraText }}</span></div></div>')

insert into WebFormQuestionTemplates([QuestionTemplateName], [QuestionTemplateText])
values ('StdTextArea',
        '<div class="row">
            <div class="col-md-12 mb-2">
                <div class="form-control">
                    <label class="form-label" for="{{ .questionID }}" style="font-size: smaller;">{{ .questionLabel }}</label>
                    <textarea class="form-control" id="{{ .questionID }}" name="{{ .InputFieldName }}" rows="{{ .questionRows }}" {{ .required }} >{{ .questionValue }}</textarea>
                </div>
            </div>
        </div>
        <div class="row"><div class="col-md-12 mb-2"><span class="fw-lighter" style="font-size: smaller;">{{ .questionExtraText }}</span></div></div>')
-----------------------------------------------------------------------------------------------------------------------------------------------------------------------

insert into [WebFormQuestions](FormID, InputFieldName, EntityName, QuestionText, QuestionType, TemplateName, Seq, AnswerRequired, QuestionExtraText)
select 4, 'BadgeName', 'BadgeName', 'Badge Name', 'text', 'stdText', 10, 'Y', 'Name as shown on your SAWG Membership Badge'

insert into [WebFormQuestions](FormID, InputFieldName, EntityName, QuestionText, QuestionType, TemplateName, Seq, AnswerRequired, QuestionExtraText)
select 4, 'BadgeNumber', 'BadgeNumber', 'Membership Number', 'number', 'stdText', 20, 'N', 'Membership # as shown on your SAWG Membership Badge'

insert into [WebFormQuestions](FormID, InputFieldName, EntityName, QuestionText, QuestionType, TemplateName, Seq, AnswerRequired)
select 4, 'Email', 'Email', 'Email', 'email', 'stdText', 30, 'N'

insert into [WebFormQuestions](FormID, InputFieldName, EntityName, QuestionText, QuestionType, TemplateName, Seq, AnswerRequired)
select 4, 'PreferredPhone', 'PreferredPhone', 'Preferred Phone or Mobile', 'tel', 'stdText', 40, 'Y'