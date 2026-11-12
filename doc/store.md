полный список новых ручек, на которые рекомендуем перейти:

Метод	Новая асинхронная ручка	Старая ручка
POST	/api/store/startCreateProject	/api/store/project
POST	/api/store/startCreateProjectFromArchive	/api/store/project/archive
POST	/api/store/project/{projectId}/startCreateBranch	/api/store/project/{projectId}/branches
POST	/api/store/project/{projectId}/startCreateBranchFromArchive	/api/store/project/{projectId}/branches/archive
POST	/api/store/uploadSession/{uploadId}/startCreateProject	/api/store/uploadSession/{uploadId}/project
POST	/api/store/uploadSession/{uploadId}/project/{projectId}/startCreateBranch	/api/store/uploadSession/{uploadId}/project/{projectId}/branches
POST	/api/store/uploadSession/{uploadId}/project/{projectId}/branches/{branchId}/startUpdate	/api/store/uploadSession/{uploadId}/project/{projectId}/branches/{branchId}
POST	/api/store/{projectId}/startUploadSources	/api/store/{projectId}/sources
POST	/api/store/{projectId}/branches/{branchId}/startUploadSources	/api/store/{projectId}/branches/{branchId}/sources
POST	/api/sourceControl/{projectId}/startCheckout	/api/sourceControl/{projectId}/checkout
POST	/api/sourceControl/{projectId}/branches/{branchId}/startUpdate	/api/sourceControl/{projectId}/branches/{branchId}/update
POST	/api/store/uploadSession/{uploadId}/cancel	/api/sourceControl/cancel/{operationId}
