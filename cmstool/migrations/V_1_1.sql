USE [subscribeproject]
GO
/****** Object:  Table [dbo].[TB_MAS_User]    Script Date: 15/12/2022 4:33:54 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[TB_TRN_Subscribers](
	[Id] [bigint] IDENTITY(1,1) NOT NULL,
	[Email] [nvarchar](255) UNIQUE NOT NULL,
	[Name] [nvarchar](255) NOT NULL,
	[IsSubscribed] [bit] NOT NULL,
    [SubscribedDate] [datetime] NOT NULL,
    [UnsubscribedDate] [datetime] NULL,
	[Delflag] [bit] NOT NULL,
 CONSTRAINT [PK_TB_TRN_Subscribers] PRIMARY KEY CLUSTERED 
(
	[Id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO
ALTER TABLE [dbo].[TB_TRN_Subscribers] ADD  CONSTRAINT [DF_TB_TRN_Subscribers_SubscribedDate]  DEFAULT (getdate()) FOR [SubscribedDate]
GO
ALTER TABLE [dbo].[TB_TRN_Subscribers] ADD  CONSTRAINT [DF_TB_TRN_Subscribers_Delflag]  DEFAULT ((0)) FOR [Delflag]
